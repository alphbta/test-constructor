import {BrowserRouter, Route, Routes} from "react-router-dom";
import Registration from "./pages/Registration";
import Home from "./pages/Home.jsx";
import ProtectedRoutes from "./utils/ProtectedRoutes.jsx"
import Login from "./pages/Login.jsx";
import Tests from "./pages/Test.jsx";
import CreateTest from "./pages/CreateTest.jsx"; // если есть

function App() {
  return (
      <BrowserRouter>
          <Routes>
              <Route path="/registration" element={<Registration />} />
              <Route path="/login" element={<Login />} />
              <Route element={<ProtectedRoutes />}>
                  <Route path="/" element={<Home />} />
                  <Route path="/tests" element={<Tests />} />
                  <Route path="/create" element={<CreateTest />} />
              </Route>
          </Routes>
      </BrowserRouter>
  )
}

export default App
