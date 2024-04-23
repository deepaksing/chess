import React, { useEffect, useState } from "react";
import axios from "axios";
import { useParams } from "react-router-dom";
import "./match.css";

function Match() {
  const matchID = parseInt(useParams().match_id, 10);
  const [chessboardState, setChessboardState] = useState([]);

  const mapNumberToPiece = (number) => {
    switch (number) {
      case -5:
        return "♜"; // Black rook
      case -4:
        return "♞"; // Black knight
      case -3:
        return "♝"; // Black bishop
      case -2:
        return "♛"; // Black queen
      case -1:
        return "♚"; // Black king
      case 1:
        return "♔"; // White king
      case 2:
        return "♕"; // White queen
      case 3:
        return "♗"; // White bishop
      case 4:
        return "♘"; // White knight
      case 5:
        return "♖"; // White rook
      default:
        return "-"; // Blank space
    }
  };

  const getBoardState = async () => {
    try {
      const response = await axios.post("http://localhost:8080/chess_board", {
        match_id: matchID,
      });
      if (response.data) {
        setChessboardState(JSON.parse(response.data));
      } else {
        console.log("No data received");
      }
    } catch (error) {
      console.error("Error fetching chessboard state:", error);
    }
  };

  useEffect(() => {
    getBoardState();
  }, []);

  // Convert 1D array to 2D array
  const chessboardRows = [];
  for (let i = 0; i < 8; i++) {
    console.log(chessboardState.slice(i * 8, (i + 1) * 8));
    chessboardRows.push(chessboardState.slice(i * 8, (i + 1) * 8));
  }

  console.log(chessboardRows);

  return (
    <div>
      <div>Chess</div>
      <div className="chessboard">
        {chessboardRows.map((row, rowIndex) => (
          <div key={rowIndex} className="chessboard-row">
            {row.map((number, columnIndex) => (
              <div key={columnIndex} className="square">
                {mapNumberToPiece(number)}
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
}

export default Match;
