import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { drawCard, resetGame } from '../actions/gameActions.js';
import { connect } from 'react-redux'; // Correct import statement

const Game = ({ count, increment, decrement }) => {
  const game = useSelector((state) => state.game);
  const dispatch = useDispatch();

  const handleDrawCard = () => {
    dispatch(drawCard());
  };

  const handleResetGame = () => {
    dispatch(resetGame());
  };

  return (
    <div>
      <h2>Game Interface</h2>
      <button onClick={handleDrawCard}>Draw Card</button>
      <button onClick={handleResetGame}>Reset Game</button>
      <p>Current Card: {game.currentCard}</p>
      <p>Player Name: {game.playerName}</p>
      <p>Points: {game.points}</p>
    </div>
  );
};

const mapStateToProps = (state) => ({
  count: state.counter.count,
});

const mapDispatchToProps = (dispatch) => ({
  increment: () => dispatch({ type: 'INCREMENT' }),
  decrement: () => dispatch({ type: 'DECREMENT' }),
});

export default connect(mapStateToProps, mapDispatchToProps)(Game);
