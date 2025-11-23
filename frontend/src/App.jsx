import { useState } from 'react'
import registerIcon from './assets/registration.svg'
import processMark from './assets/process-mark.svg'
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Registration from "./pages/Registration";
import Home from "./pages/Home.jsx";
import ProtectedRoutes from "./utils/ProtectedRoutes.jsx"
import Login from "./pages/Login.jsx";

function App() {
  return (
      <BrowserRouter>
          <Routes>
              <Route path="/registration" element={<Registration />} />
              <Route path="/login" element={<Login />} />
              <Route element={<ProtectedRoutes />}>
                  <Route path="/" element={<Home />} />
              </Route>
          </Routes>
      </BrowserRouter>
  )
}

export default App
