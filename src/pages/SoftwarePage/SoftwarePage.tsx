import { type FC, useState, useEffect } from "react";
import { Row, Col, Spinner, Container } from "react-bootstrap";
import InputField from "../../components/InputField/InputField";
import { SoftwareCard } from "../../components/SoftwareCard/SoftwareCard";
import { CartFloating } from "../../components/CartFloating/CartFloating";
import { getSoftwareByName, type ISoftware } from "../../modules/softwareApi";
import { SOFTWARE_MOCK } from "../../modules/mock";
import "./SoftwarePage.css";
import { useNavigate } from "react-router-dom";
import { BreadCrumbs } from "../../components/BreadCrumbs/BreadCrumbs";

export const SoftwarePage: FC = () => {
    const [searchValue, setSearchValue] = useState("");
    const [loading, setLoading] = useState(false);
    const [software, setSoftware] = useState<ISoftware[]>([]);
    const [cartCount, setCartCount] = useState(0);
    const navigate = useNavigate();

    // Определяем крошки для этой страницы
    const breadcrumbs = [
        { label: "Программное обеспечение", path: undefined } // текущая страница, не кликабельна
    ];

    const fetchData = async (name = "") => {
        setLoading(true);
        
        try {
            const data = await getSoftwareByName(name);
            
            if (!data || !Array.isArray(data)) {
                throw new Error("Invalid response format");
            }
            
            setSoftware(data);
        } catch (error) {
            // Используем mock данные при ошибке без показа сообщений
            const filteredMock = SOFTWARE_MOCK.filter(item =>
                item.title.toLowerCase().includes(name.toLowerCase())
            );
            
            setSoftware(filteredMock);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData("");
    }, []);

    const handleSearch = () => {
        fetchData(searchValue);
    };

    const handleCardClick = (id: number) => {
        navigate(`/software/${id}`);
    };

    const handleAdd = (id: number) => {
        console.log("Добавлено в корзину: ", id);
        setCartCount(prev => prev + 1);
    };

    const handleCartClick = () => {
        if (cartCount > 0) {
            navigate("/cart");
        }
    };

    return (
        <>
            {/* BreadCrumbs теперь вне контейнера */}
            <BreadCrumbs crumbs={breadcrumbs} />
            
            <Container className="page-container">
                <InputField
                    value={searchValue}
                    setValue={setSearchValue}
                    onSubmit={handleSearch}
                    placeholder="Быстрый поиск"
                    loading={loading}
                />

                <h2 className="software-title-center">Серверное ПО</h2>

                {loading && (
                    <div className="loadingBg">
                        <Spinner animation="border" />
                    </div>
                )}

                {!loading && (
                    <>
                        {software.length === 0 ? (
                            <h3 className="empty-msg">Ничего не найдено</h3>
                        ) : (
                            <Row xs={1} md={2} lg={3} className="g-4 cards-grid">
                                {software.map((item) => (
                                    <Col key={item.software_id}>
                                        <SoftwareCard
                                            img_url={item.img_url}
                                            title={item.title}
                                            version={item.version}
                                            size={item.size}
                                            description={item.description}
                                            onOpen={() => handleCardClick(item.software_id)}
                                            onAdd={() => handleAdd(item.software_id)}
                                        />
                                    </Col>
                                ))}
                            </Row>
                        )}
                    </>
                )}

                <CartFloating count={cartCount} onOpen={handleCartClick} />
            </Container>
        </>
    );
};