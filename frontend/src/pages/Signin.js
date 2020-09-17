import React from 'react';
import { Button, Form } from 'react-bootstrap';
import { request } from '../common/utils.js'

let state = {
  username: '',
  password: ''
}

function Signin(props) {  
  return (
    <>
      <Form onSubmit={ send }>
        <Form.Group controlId="formUsername">
          <Form.Label>Username</Form.Label>
          <Form.Control 
            type="text"
            name="username"
            placeholder="Username" 
            onChange={handleChange}
          />
          <Form.Text className="text-muted">
          </Form.Text>
        </Form.Group>
        <Form.Group controlId="formPassword">
          <Form.Label>Password</Form.Label>
          <Form.Control 
            type="password" 
            name="password"
            placeholder="Password" 
            onChange={ handleChange }
          />
        </Form.Group>
        <Button variant="primary" type="submit">
          Sign In
        </Button>
      </Form>
    </>
  );
}

function send(e) {
  e.preventDefault();
  console.log(e);
  request('POST', 'http://localhost:80/api/auth/signin', {}, JSON.stringify(state))
    .then((res) => {
      
    })
    .catch((e) => {
      console.log("err: " + e);
    });
}

function handleChange(e) {
  state[e.target.name] = e.target.value;
}

export default Signin;
