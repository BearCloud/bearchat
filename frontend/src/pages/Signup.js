import React, { useState } from 'react';
import { Button, Form } from 'react-bootstrap';
import { request } from '../common/utils.js'

function Signup(props) {  

  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('' );
  
  const send = (e) => {
    e.preventDefault();
    request('POST', 'http://localhost:80/api/auth/signup', {}, JSON.stringify({ email, username, password }))
      .then((res) => {
        console.log(res.status)
      })
      .catch((e) => {
        console.log("err: " + e);
      });
  }

  return (
    <>
      <Form onSubmit={ send }>
        <Form.Group controlId="formBasicEmail">
          <Form.Label>Email address</Form.Label>
          <Form.Control 
            type="email" 
            name="email"
            placeholder="Enter email"
            onChange={(e) => setEmail(e.target.value)}
          />
          <Form.Text className="text-muted">
            We'll never share your email with anyone else.
          </Form.Text>
        </Form.Group>
        <Form.Group controlId="formUsername">
          <Form.Label>Username</Form.Label>
          <Form.Control 
            type="text"
            name="username"
            placeholder="Username" 
            onChange={(e) => setUsername(e.target.value)}
          />
        </Form.Group>
        <Form.Group controlId="formPassword">
          <Form.Label>Password</Form.Label>
          <Form.Control 
            type="password" 
            name="password"
            placeholder="Password" 
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Group>
        <Button variant="primary" type="submit">
          Sign Up
        </Button>
      </Form>
    </>
  );
}

export default Signup;
