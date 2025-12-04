import { Navbar, Nav } from "react-bootstrap";
import { Link } from "react-router-dom";
import "./Navigation.css";

const Navigation = () => {
  return (
    <Navbar bg="light" expand="lg" className="top-navbar no-container-navbar">
      <Navbar.Brand as={Link} to="/" className="navbar-logo">
        IT Soft
      </Navbar.Brand>
      
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="ms-auto">
          <Nav.Link as={Link} to="/software" className="services-link">
            Услуги
          </Nav.Link>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
};

export default Navigation;