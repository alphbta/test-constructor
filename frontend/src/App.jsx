import {BrowserRouter, Route, Routes, Navigate} from "react-router-dom";
import Registration from "./pages/Registration";
import Home from "./pages/Home.jsx";
import ProtectedRoutes from "./utils/ProtectedRoutes.jsx"
import Login from "./pages/Login.jsx";
import Tests from "./pages/Test.jsx";
import CreateTest from "./pages/CreateTest.jsx";

import MyTestStudent from "./pages/MyTestStudent.jsx";
import StatisticsTest from "./pages/StatisticsTest.jsx";
import PassingTestStudent from "./pages/PassingTestStudent.jsx";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navigate to="/registration" replace />} />

                <Route path="/registration" element={<Registration />} />
                <Route path="/login" element={<Login />} />

                <Route path="/myTestStudent" element={<MyTestStudent />} />
                <Route path="/statisticsTest" element={<StatisticsTest />} />
                <Route path="/passingTestStudent" element={<PassingTestStudent />} />


                <Route path="/statistics/:testId" element={<StatisticsTest />} />
                <Route path="/test/:test_link" element={<PassingTestStudent />} />

                <Route element={<ProtectedRoutes />}>
                    <Route path="/home" element={<Home />} />
                    <Route path="/tests" element={<Tests />} />
                    <Route path="/create" element={<CreateTest />} />
                </Route>
            </Routes>
        </BrowserRouter>
    )
}

export default App