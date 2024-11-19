import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from "./Pages/home.js";
import Leaderboard from './Components/leaderboard.js';
import Game from "./Pages/game.js"

function App() {
  return (
    <Router>
      <h1>Exploding Kitten Game</h1>
      {/* <Game />
      <Leaderboard /> */}
      <Routes>
        <Route path="/" element={<Home />} />        
        <Route path="/game" element={<Game />} />
      </Routes>
    </Router>
  );
}

export default App;
