import { type FC, useEffect } from "react";
import { Container, Row, Col } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import startGif from "../../assets/start.gif";
import "./HomePage.css";

export const HomePage: FC = () => {
    const navigate = useNavigate();

    useEffect(() => {
        document.body.classList.add('no-scroll');
        
        return () => {
            document.body.classList.remove('no-scroll');
        };
    }, []);

    const handleGifClick = () => {
        navigate("/software");
    };

    return (
        <div className="home-page">
            <Container className="home-container">
                <Row className="justify-content-center">
                    <Col xs={12} className="text-center">
                        {/* Анимированный GIF с обработчиком клика */}
                        <div className="animation-container">
                            <img 
                                src={startGif} 
                                alt="IT Soft Animation" 
                                className="start-animation clickable-gif"
                                onClick={handleGifClick}
                            />
                        </div>
                        
                        {/* Основной заголовок */}
                        <div className="main-title-container">
                            <h1 className="main-title">
                                IT Soft
                            </h1>
                            <p className="subtitle">
                                Cервис по установке серверного программного обеспечения
                            </p>
                        </div>
                    </Col>
                </Row>
            </Container>
        </div>
    );
};