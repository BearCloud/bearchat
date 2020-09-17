import React, { useState }  from 'react';
import { Button, Form } from 'react-bootstrap';
import { request } from '../common/utils.js'

function Signin(props) {  
  
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('' );

  const send = (e) => {
    e.preventDefault();
    request('POST', 'http://localhost:80/api/auth/signin', {}, JSON.stringify({ username, password }))
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
          Sign In
        </Button>
      </Form>
    </>
  );
}

export default Signin;
