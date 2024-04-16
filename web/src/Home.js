import React, { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

function Home() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
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
        navigate("/game");
      } else {
        // Handle unsuccessful login (e.g., display error message)
        setError("Login failed. Please try again.");
      }
    } catch (error) {
      setError(error.response.data.error);
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
