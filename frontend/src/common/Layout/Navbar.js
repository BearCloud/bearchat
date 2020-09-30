import React, { useState } from 'react';
import ReactNav from 'react-bootstrap/Nav';
import ReactNavbar from 'react-bootstrap/Navbar';
import './Navbar.css';
import { request, getUUID } from '../utils.js';

function Navbar(props) {
  const [isAuth, setIsAuth] = useState(false);
  const uuid = getUUID();
  request('GET', 'http://localhost:81/api/posts/0', {})
    .then((res) => {
      console.log(res.status, res.responseText);
    })
    .catch((res) => {
      console.log("ping err: ", res);
    });

  return (
    <ReactNavbar bg="light" variant="light">
      <ReactNavbar.Brand href="/">BearChat</ReactNavbar.Brand>
      <ReactNav className="mr-auto">
        <ReactNav.Link href="/signin">Sign In</ReactNav.Link>
        <ReactNav.Link href="/signup">Sign Up</ReactNav.Link>
      </ReactNav>
    </ReactNavbar>
  );
}

export default Navbar;
