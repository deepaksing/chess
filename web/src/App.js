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
import Match from "./Match";

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
  {
    path: "/match",
    element: <Match />,
  },
]);

// App component
const App = () => (
  <RouterProvider router={router}>
    <div className="App">
      <h1>Welcome to the React Router App</h1>
      <Route path="/" element={<Home />} />
      <Route path="/game" element={<Game />} />
      <Route path="/match" element={<Match />} />
    </div>
  </RouterProvider>
);

export default App;
