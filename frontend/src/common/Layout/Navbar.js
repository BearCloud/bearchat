import React from 'react';
import ReactNav from 'react-bootstrap/Nav';
import ReactNavbar from 'react-bootstrap/Navbar';
import './Navbar.css';

function Navbar(props) {
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
