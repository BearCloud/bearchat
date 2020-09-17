import decode from 'jsonwebtoken'
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
  let decoded = decode(token, {complete: true});
  return decoded.payload.UserID;
}

export function request(method, url, qs, body) {
  const allowedOrigins = ["http://localhost:3000", "http://localhost:80"];
  return new Promise((resolve, reject) => {
    let xhr = new XMLHttpRequest();
    let u = new URL(url);
    for (const [key, value] of Object.entries(qs)) {
      u.searchParams.append(key, value);
    }
    xhr.open(method, u.toString(), true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Content-Length", body.length);
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve(xhr);
      } else {
        reject(xhr.status);
      }
    };
    xhr.onerror = () => {
      reject(xhr.status);
    }
    xhr.send(body);
  });
}
