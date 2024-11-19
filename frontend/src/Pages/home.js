import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const Home = () => {
  const [username, setUsername] = useState("");
  const navigate=useNavigate()

//   useEffect(() => {
//     const fetchLeaderboard = async () => {
//       const response = await axios.get("http://localhost:8000/api/leaderboard");
//       setLeaderboard(response.data);
//     };

//     fetchLeaderboard();
//   }, []);
// useEffect(() => {
//   return () => {
//     if (websocket) {
//       websocket.close();
//     }
//   };
// }, []);


  const handleSubmit = async (e) => {
    e.preventDefault();

    const data={
        username:username,
        points:0
    }
    console.log(data)
    await axios.post("http://localhost:8000/create", data)
    .then(response => {
      console.log(response)
      // Collect cookies from the response headers
      const cookies = response.headers['set-cookie'];
      console.log(cookies);
  
      // Extract and set the cookies in the document
      cookies.forEach(cookie => {
        document.cookie = cookie.split(';')[0];
      });
  
      // Continue with your logic...
    })
    .catch(error => {
      console.error('Error:', error);
    });


  //   const socket = new WebSocket('ws://localhost:8000/ws');

  //   socket.addEventListener('open', (event) => {
  //     socket.send(JSON.stringify(data)); // Send the username as a JSON string
  // });
  //   socket.addEventListener('message', (event) => {
  //     const token = event.data;
  //     console.log('Received token:', token);
  
  //     // Set the token in a cookie
  //     document.cookie = `token=${token}; path=/`;
  // });
  //   setWebsocket(socket);
    navigate("/game")
  };

  const handleD=async(e)=>{
    e.preventDefault();

    const data={
      username: "madhav",
      password: "password"
    }

    const res=await axios.post("http://localhost:8000/signin",data)
    console.log(res)
  }

  return (
    <div>
      <h1>Exploding Kitten</h1>
      <div>
        <form onSubmit={handleSubmit}>
            <input type="text" onChange={(e)=>setUsername(e.target.value)} required/>
            <button type="submit">Submit</button>
        </form>
        <button onClick={handleD}>teting</button>
      </div>

    </div>
  );
};

export default Home;
