import React, { useState, useEffect } from "react";
import axios from "axios";

const Game = () => {
  const [websocket, setWebsocket] = useState(null);
  const [token, setToken] = useState("");
  const [cards, setCards] = useState([]);

  // Function to fetch and set random cards
  const fetchRandomCards = async () => {
    try {
      const response = await axios.get('ws://localhost:8000/fetch');
      console.log(response)
      setCards(response.data);
    } catch (error) {
      console.error('Failed to fetch random cards:', error);
    }
  };

  useEffect(() => {
    console.log(document.cookie.split(';'))
    // const token = getCookie("token");
    // if (token) {
    //   setToken(token);
    //   fetchRandomCards()
    // }
});

useEffect(() => {
  return () => {
    if (websocket) {
      websocket.close();
    }
  };
}, []);

  const handleClick = async (e) => {
    e.preventDefault();

    const data={
        token:token,
        points:0
    }
    console.log(data)
    const socket = new WebSocket('ws://localhost:8000/game');

    socket.addEventListener('open', (event) => {
      socket.send(JSON.stringify(data)); // Send the username as a JSON string
  });

    socket.addEventListener('message', (event) => {
        console.log('Message from server: ', event.data);
    });

    setWebsocket(socket);
};

  return (
    <div>
      <h1>Play</h1>
      <div>
      {cards.length > 0 && cards.map((card, index) => (
          <div key={index}>{card.type}</div>
        ))}
      </div>
      <div>
        <button onClick={handleClick}>Cat</button>
        <button onClick={handleClick}>Defuse</button>
        <button onClick={handleClick}>Shuffle</button>
        <button onClick={handleClick}>Exploding cat</button>
      </div>

    </div>
  );
};

export default Game;
