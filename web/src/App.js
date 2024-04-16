import * as React from "react";
import { createRoot } from "react-dom/client";
import {
  createBrowserRouter,
  RouterProvider,
  Route,
  Link,
} from "react-router-dom";
import Home from "./Home";
import Game from "./Game";

// Create browser router instance
const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/game",
    element: <Game />,
  },
]);

// App component
const App = () => (
  <RouterProvider router={router}>
    <div className="App">
      <h1>Welcome to the React Router App</h1>
      <Route path="/" element={<Home />} />
      <Route path="/game" element={<Game />} />
    </div>
  </RouterProvider>
);

export default App;
