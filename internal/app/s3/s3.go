package s3

import (
    "context"
    "fmt"
    "mime/multipart"
    "path/filepath"
    "strings"
    "time"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "github.com/sirupsen/logrus"
    "lab1/internal/app/config"
    "lab1/internal/app/repository"
)

//
// --- ИНИЦИАЛИЗАЦИЯ MinIO ---
//
func InitMinioClient(conf *config.MinioConfig) (*minio.Client, error) {
    client, err := minio.New(conf.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
        Secure: conf.UseSSL,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to init MinIO client: %w", err)
    }

    ctx := context.Background()

    _, err = client.ListBuckets(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
    }

    exists, err := client.BucketExists(ctx, conf.BucketName)
    if err != nil {
        return nil, fmt.Errorf("failed to check bucket existence: %w", err)
    }

    if !exists {
        logrus.Infof("Bucket %s does not exist, creating...", conf.BucketName)
        err = client.MakeBucket(ctx, conf.BucketName, minio.MakeBucketOptions{})
        if err != nil {
            return nil, fmt.Errorf("failed to create bucket: %w", err)
        }

        policy := fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": "*"},
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::%s/*"]
                }
            ]
        }`, conf.BucketName)

        err = client.SetBucketPolicy(ctx, conf.BucketName, policy)
        if err != nil {
            logrus.Warnf("Failed to set bucket policy: %v", err)
        }
    }

    logrus.Info("Successfully connected to MinIO and ensured bucket exists")
    return client, nil
}

//
// --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ---
//

// extractFileNameFromURL достаёт имя файла из MinIO URL
func extractFileNameFromURL(url string) string {
    parts := strings.Split(url, "/")
    if len(parts) == 0 {
        return ""
    }
    return parts[len(parts)-1]
}

//
// --- ЗАГРУЗКА ИЗОБРАЖЕНИЯ ---
//
func UploadSoftwareImage(
    ctx context.Context,
    repo *repository.Repository,
    minioClient *minio.Client,
    conf *config.MinioConfig,
    softwareID uint,
    fileHeader *multipart.FileHeader,
) (string, int64, string, error) {

    software, err := repo.GetSoftwareByID(softwareID)
    if err != nil {
        return "", 0, "", fmt.Errorf("software not found: %w", err)
    }

    if fileHeader.Size > 10*1024*1024 {
        return "", 0, "", fmt.Errorf("file size too large (max 10MB)")
    }

    ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
    allowed := map[string]string{
        ".jpg":  "image/jpeg",
        ".jpeg": "image/jpeg",
        ".png":  "image/png",
        ".gif":  "image/gif",
        ".webp": "image/webp",
    }
    contentType, ok := allowed[ext]
    if !ok {
        return "", 0, "", fmt.Errorf("invalid file type: %s", ext)
    }

    file, err := fileHeader.Open()
    if err != nil {
        return "", 0, "", fmt.Errorf("failed to open uploaded file: %w", err)
    }
    defer file.Close()

    bucket := conf.BucketName
    newFileName := fmt.Sprintf("software_%d_%d%s", softwareID, time.Now().Unix(), ext)

    // Удаляем старый файл, если есть
    if software.ImgURL != "" {
        oldFileName := extractFileNameFromURL(software.ImgURL)
        if oldFileName != "" {
            err = minioClient.RemoveObject(ctx, bucket, oldFileName, minio.RemoveObjectOptions{})
            if err != nil {
                logrus.Warnf("Failed to remove old image: %v", err)
            }
        }
    }

    // Загружаем новый
    info, err := minioClient.PutObject(ctx, bucket, newFileName, file, fileHeader.Size, minio.PutObjectOptions{
        ContentType: contentType,
        UserMetadata: map[string]string{
            "software-id": fmt.Sprint(softwareID),
            "upload-time": time.Now().Format(time.RFC3339),
        },
    })
    if err != nil {
        return "", 0, "", fmt.Errorf("failed to upload image: %w", err)
    }

    imgURL := fmt.Sprintf("http://%s/%s/%s", conf.Endpoint, bucket, newFileName)

    if err := repo.UpdateSoftwareImage(softwareID, imgURL); err != nil {
        _ = minioClient.RemoveObject(ctx, bucket, newFileName, minio.RemoveObjectOptions{})
        return "", 0, "", fmt.Errorf("failed to update database: %w", err)
    }

    return imgURL, info.Size, newFileName, nil
}

//
// --- УДАЛЕНИЕ ИЗОБРАЖЕНИЯ ---
//
func DeleteSoftwareImage(
    ctx context.Context,
    minioClient *minio.Client,
    conf *config.MinioConfig,
    imageURL string,
) error {
    if imageURL == "" {
        return nil
    }
    fileName := extractFileNameFromURL(imageURL)
    if fileName == "" {
        return fmt.Errorf("invalid image URL: %s", imageURL)
    }

    err := minioClient.RemoveObject(ctx, conf.BucketName, fileName, minio.RemoveObjectOptions{})
    if err != nil {
        return fmt.Errorf("failed to remove object from MinIO: %w", err)
    }

    logrus.Infof("Removed image from MinIO: %s", fileName)
    return nil
}
