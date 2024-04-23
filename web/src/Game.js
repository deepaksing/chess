import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function Game() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const fetchIntervalRef = useRef(0);
  const [opponent, setOpponent] = useState("opponent");
  const [matchID, setMatchID] = useState(null);

  // useEffect(() => {
  //   clearInterval(fetchStatusInterval);
  // }, [opponent]);

  const getAnotherPlayer = async (username) => {
    try {
      const response = await axios.post("http://localhost:8080/join", {
        username: username,
      });
      if (response) {
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
        const isPlyaing = response.data[0].isplaying;
        if (isPlyaing) {
          setOpponent(response.data[0].opponent);
          setMatchID(response.data[0].match_id);
          clearInterval(fetchIntervalRef.current);
          navigate(`/match/${response.data[0].match_id}`);
        }
      } else {
      }
    } catch (error) {
      setError(error.response.data.error);
    }
  };

  console.log("opponent ", opponent);
  console.log("mfetch interval 2 ", fetchIntervalRef.current);

  const joinGame = () => {
    const username = localStorage.getItem("username");
    if (username) {
      setLoading(true);
      getAnotherPlayer(username);
      var intervalId = setInterval(() => {
        getStatus(username);
      }, 5000);
      // console.log(intervalId);
      fetchIntervalRef.current = intervalId;
      console.log("Fetch Interval ID:", intervalId);
    }
  };

  return (
    <div className="game">
      {loading ? <div>Loading</div> : null}
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
