import React, { useState, useEffect } from "react";
import axios from "axios";

const Leaderboard = () => {
  const [leaderboard, setLeaderboard] = useState({});

  useEffect(() => {
    const fetchLeaderboard = async () => {
      const response = await axios.get("http://localhost:8000/api/leaderboard");
      setLeaderboard(response.data);
    };

    fetchLeaderboard();
  }, []);

  const updatePoints = async (username, points) => {
    await axios.post("http://localhost:8000/api/update", { username, points });
    const response = await axios.get("http://localhost:8000/api/leaderboard");
    setLeaderboard(response.data);
  };

  return (
    <div>
      <h1>Leaderboard</h1>
      <table>
        <thead>
          <tr>
            <th>Username</th>
            <th>Points</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(leaderboard).map(([username, points]) => (
            <tr key={username}>
              <td>{username}</td>
              <td>{points}</td>
              <td>
                <button onClick={() => updatePoints(username, 1)}>Add Point</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Leaderboard;
