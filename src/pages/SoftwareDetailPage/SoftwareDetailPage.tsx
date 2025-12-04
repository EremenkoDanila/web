import { type FC, useState, useEffect } from "react";
import { Container, Row, Col, Spinner } from "react-bootstrap";
import { useParams, useNavigate } from "react-router-dom";
import { CartFloating } from "../../components/CartFloating/CartFloating";
import { getSoftwareById, type ISoftware } from "../../modules/softwareApi";
import { SOFTWARE_MOCK } from "../../modules/mock";
import defaultImage from "../../assets/program.png";
import "./SoftwareDetailPage.css";
import { BreadCrumbs } from "../../components/BreadCrumbs/BreadCrumbs";

export const SoftwareDetailPage: FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [software, setSoftware] = useState<ISoftware | null>(null);
    const [loading, setLoading] = useState(true);
    const [cartCount, setCartCount] = useState(0);
    const [imageError, setImageError] = useState(false);
    const [breadcrumbs, setBreadcrumbs] = useState([
        { label: "Программное обеспечение", path: "/software" },
        { label: "Загрузка...", path: undefined }
    ]);

    // Убираем скролл при монтировании
    useEffect(() => {
        document.body.classList.add("no-scroll");
        return () => document.body.classList.remove("no-scroll");
    }, []);

    useEffect(() => {
        const fetchSoftwareDetail = async () => {
            if (!id) return;

            setLoading(true);
            try {
                const data = await getSoftwareById(parseInt(id));

                // Проверяем корректность данных
                if (!data || typeof data !== "object") {
                    throw new Error("Invalid API response");
                }

                setSoftware(data);
                setBreadcrumbs([
                    { label: "Программное обеспечение", path: "/software" },
                    { label: data.title, path: undefined }
                ]);
            } catch (error) {
                console.error("Error fetching software details:", error);

                // Подключаем mock данные
                const fallback = SOFTWARE_MOCK.find(
                    (item) => item.software_id === parseInt(id)
                );

                if (fallback) {
                    setSoftware(fallback);
                    setBreadcrumbs([
                        { label: "Программное обеспечение", path: "/software" },
                        { label: fallback.title, path: undefined }
                    ]);
                } else {
                    setSoftware(null);
                    setBreadcrumbs([
                        { label: "Программное обеспечение", path: "/software" },
                        { label: "Не найдено", path: undefined }
                    ]);
                }
            } finally {
                setLoading(false);
            }
        };

        fetchSoftwareDetail();
    }, [id]);

    const handleCartClick = () => {
        if (cartCount > 0) {
            console.log("Navigate to cart");
        }
    };

    const handleImageError = () => {
        console.warn("Failed to load image:", software?.img_url);
        setImageError(true);
    };

    const imageSrc =
        !software?.img_url || imageError ? defaultImage : software.img_url;

    if (loading) {
        return (
            <>
                <BreadCrumbs crumbs={breadcrumbs} />
                <div className="software-detail-page">
                    <Container className="detail-container">
                        <div className="loading-container">
                            <Spinner animation="border" variant="success" />
                        </div>
                    </Container>
                </div>
            </>
        );
    }

    if (!software) {
        return (
            <>
                <BreadCrumbs crumbs={breadcrumbs} />
                <div className="software-detail-page">
                    <Container className="detail-container">
                        <div className="error-container">
                            <h2>Программное обеспечение не найдено</h2>
                            <button onClick={() => navigate("/software")}>
                                Вернуться к списку
                            </button>
                        </div>
                    </Container>
                </div>
            </>
        );
    }

    return (
        <>
            <BreadCrumbs crumbs={breadcrumbs} />
            
            <div className="software-detail-page">
                <Container className="detail-container">
                    <Row className="product-layout">
                        {/* Изображение */}
                        <Col md={4} className="image-section">
                            <img
                                src={imageSrc}
                                alt={`${software.title} Logo`}
                                className="product-image"
                                onError={handleImageError}
                            />
                        </Col>

                        {/* Контент */}
                        <Col md={8} className="content-section">
                            <h1 className="product-title">{software.title}</h1>

                            <section className="support-info">
                                <div className="support-row">
                                    <span className="label">Поддержка:</span>
                                    <span className="value">{software.os}</span>
                                </div>
                                <div className="support-row">
                                    <span className="label">Версия:</span>
                                    <span className="value">{software.version}</span>
                                </div>
                                <div className="support-row">
                                    <span className="label">Размер:</span>
                                    <span className="value">{software.size} ГБ</span>
                                </div>
                            </section>
                        </Col>
                    </Row>

                    {/* Описание */}
                    <section className="description">
                        <p>{software.full_description}</p>
                    </section>

                    {/* Функции */}
                    <section className="functions">
                        <h2>Основные функции:</h2>
                        <div className="functions-text">
                            {software.functions.split("\\n").map((func, index) => (
                                <div key={index} className="function-item">
                                    {func}
                                </div>
                            ))}
                        </div>
                    </section>
                </Container>

                <CartFloating count={cartCount} onOpen={handleCartClick} />
            </div>
        </>
    );
};