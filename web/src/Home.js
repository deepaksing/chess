import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

function Home() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [receivedMessage, setReceivedMessage] = useState("");
  const [ws, setWs] = useState(null);
  const navigate = useNavigate(); // Hook to navigate programmatically

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      // Make a POST request to login endpoint
      const response = await axios.post("http://localhost:8080/login", {
        username: username,
        password: password,
      });
      // Assuming the response contains a success field indicating successful login
      if (response) {
        // Navigate to the /game route upon successful login
        localStorage.setItem("username", username);
        navigate("/game");
      } else {
        // Handle unsuccessful login (e.g., display error message)
        setError("Login failed. Please try again.");
      }
    } catch (error) {
      setError(error.response.data.error);
    }
  };

  useEffect(() => {
    // Connect to WebSocket server
    const socket = new WebSocket("ws://localhost:8080/match");
    socket.onopen = () => {
      console.log("WebSocket connected");
    };

    socket.onmessage = (event) => {
      setReceivedMessage(event.data);
    };

    setWs(socket);

    return () => {
      socket.close();
    };
  }, []); // Empty dependency array ensures this effect runs only once

  const sendMessage = () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(message);
    }
  };
  return (
    <div className="main">
      <div>
        <form className="login_form" onSubmit={handleLogin}>
          <input
            type="text"
            placeholder="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="password"
            placeholder="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <button type="submit" className="login_btn">
            Join
          </button>
        </form>
        {error && <div className="error">{error}</div>}
      </div>
    </div>
  );
}

export default Home;
