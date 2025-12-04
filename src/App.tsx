import { BrowserRouter, Routes, Route } from "react-router-dom";
import Navigation from "./components/Navigations/Navigations";
import { SoftwarePage } from "./pages/SoftwarePage/SoftwarePage";
import { HomePage } from "./pages/HomePage/HomePage";
import { SoftwareDetailPage } from "./pages/SoftwareDetailPage/SoftwareDetailPage"; // Раскомментируем

function App() {
  return (
    <BrowserRouter>
      <Navigation />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/software" element={<SoftwarePage />} />
        <Route path="/software/:id" element={<SoftwareDetailPage />} /> {/* Раскомментируем */}
      </Routes>
    </BrowserRouter>
  );
}

export default App;