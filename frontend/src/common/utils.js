import React from 'react';
import decode from 'jwt-decode'
import Cookies from 'universal-cookie';

const cookies = new Cookies();

export function getUsername() {
  let loginToken = cookies.get("access_token");
  let uuid = getUUIDFromToken(loginToken);
  // TODO: convert to username
  return uuid;
}

function getUUIDFromToken(token) {
  // TODO: verify signature; return error if invalid
  // TODO: verify header and payload; return error if invalid
  let decoded = jwt.decode(token, {complete: true});
  return decoded.payload.UserID;
}

