import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function Game() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [fetchStatusInterval, setFetchStatusInterval] = useState(0);
  const [opponent, setOpponent] = useState("opponent");
  const [matchID, setMatchID] = useState(null);

  useEffect(() => {
    console.log(opponent);
    console.log(fetchStatusInterval);
    clearInterval(fetchStatusInterval);
  }, [opponent]);

  const getAnotherPlayer = async (username) => {
    try {
      const response = await axios.post("http://localhost:8080/join", {
        username: username,
      });
      if (response) {
        console.log(response);
        // navigate("/");
      } else {
      }
    } catch (error) {
      setError(error.response.data.error);
    }
  };

  const getStatus = async (username) => {
    try {
      const response = await axios.post("http://localhost:8080/get-status", {
        username: username,
      });
      // Assuming the response contains a success field indicating successful login
      if (response) {
        setLoading(false);
        console.log(response);
        const isPlyaing = response.data[0].isplaying;
        console.log(isPlyaing);
        if (isPlyaing) {
          console.log(
            response.data[0].opponent + " " + response.data[0].match_id
          );
          setOpponent(response.data[0].opponent);
          setMatchID(response.data[0].match_id);
          // navigate("/match");
        }
      } else {
      }
    } catch (error) {
      setError(error.response.data.error);
    }
  };

  const joinGame = () => {
    const username = localStorage.getItem("username");
    console.log(username);
    if (username) {
      setLoading(true);
      getAnotherPlayer(username);
      const intervalId = setInterval(() => {
        getStatus(username);
      }, 5000);
      // console.log(intervalId);
      // setFetchStatusInterval(intervalId);
    }
  };

  return (
    <div className="game">
      {loading ? (
        <div>Loading</div>
      ) : (
        <div>
          {opponent} {fetchStatusInterval}
        </div>
      )}
      <h1>Chess Game</h1>
      <div>
        <div onClick={joinGame}>Join Game</div>
        <div>History</div>
        {error}
      </div>
      {/* Add your game components here */}
    </div>
  );
}

export default Game;
