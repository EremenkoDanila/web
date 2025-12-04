import { type FC, useState } from "react";
import { Card, Button } from "react-bootstrap";
import defaultImage from "../../assets/program.png";
import "./SoftwareCard.css";

export interface ISoftwareCardProps {
    img_url: string;
    title: string;
    version: string;
    size: number;
    description: string;
    onOpen: () => void;
    onAdd: () => void;
}

export const SoftwareCard: FC<ISoftwareCardProps> = ({
    img_url,
    title,
    version,
    size,
    description,
    onOpen,
    onAdd
}) => {
    const [imageError, setImageError] = useState(false);

    const handleImageError = () => {
        console.warn("Failed to load image:", img_url);
        setImageError(true);
    };

    // Если img_url пустая или не определена, сразу используем дефолтную картинку
    const imageSrc = !img_url || imageError ? defaultImage : img_url;

    return (
        <Card className="software-card">
            <div className="card-header-area" onClick={onOpen}>
                <Card.Title className="software-title">{title}</Card.Title>
            </div>

            <Card.Body className="card-body">
                <Card.Text className="software-info">
                    Версия: {version}
                </Card.Text>
                <Card.Text className="software-info">
                    Размер: {size} ГБ
                </Card.Text>

                <Card.Text className="software-desc">
                    {description}
                </Card.Text>

                <Card.Img
                    variant="top"
                    src={imageSrc}
                    className="software-img"
                    onClick={onOpen}
                    onError={handleImageError}
                    alt={title}
                />

                {/* <Button className="add-btn" onClick={onAdd}>
                    <span className="btn-text">Добавить в заявку</span>
                    <span className="btn-icon">+</span>
                </Button> */}
            </Card.Body>
        </Card>
    );
};