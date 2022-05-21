import React from "react";
import "./App.css";

import { BrowserRouter, Route, Routes } from "react-router-dom";

import Top from "./components/pages/Top";
import Room from "./components/pages/Room";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Top />} />
        <Route path="/rooms/:id" element={<Room />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
