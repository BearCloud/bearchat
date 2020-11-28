import React, { useState }  from 'react';
import { Button, Form, Card, InputGroup, FormControl } from 'react-bootstrap';
import { request, getUUID } from '../common/utils.js';
import swal from 'sweetalert';

import { useParams } from "react-router-dom";

function Profile(props) {

  let { uuid } = useParams();

  const ourUUID = getUUID();

  console.log("Requested profile ID:", uuid);
  console.log("Our uuid:", ourUUID);

  const [profile, setProfile] = useState(null);

  if (profile === null) {
    request('GET', `http://localhost:82/api/profile/${ourUUID}`, {})
        .then((res) => {
          // console.log(res.responseText);
          setProfile(JSON.parse(res.responseText));
        })
        .catch(() => {
          console.error("Could not retrieve profile!");
        })
    ;
  }

  const send = (e) => {
    e.preventDefault();

    const content = {
      "firstName": e.target[0].value,
      "lastName": e.target[1].value,
      "uuid": ourUUID,
      "email": e.target[2].value
    };

    console.log("Profile formContent:", content);

    request('PUT', `http://localhost:82/api/profile/${ourUUID}`, {}, JSON.stringify({ content }))
      .then((res) => {
        console.log(res.status);
        swal({
          title: "Updated!",
          text: "Successfully updated profile!",
          icon: "success",
          timeout: 5000
        }).then(() => {
          window.location.reload();
        });
      })
      .catch((res) => {
        console.log("err: ", res);
        const errMessage = res?.responseText?.trim();
        swal({
          title: "Could not update profile!",
          text: `Error when attempting to update profile (HTTP ${res.status}): ${errMessage}.`,
          icon: "error"
        });
      });
  };

  const thisIsUs = !uuid || ourUUID === uuid; // if we are looking at ourselves or not

  var profileHtml = [];
  if (profile) {
    profileHtml = (
      <Card style={{ width: '35rem' }}>
        <Card.Body>
          <Card.Title>{profile.firstName} {profile.lastName}</Card.Title>
          <Card.Subtitle className="mb-2 text-muted">User ID {profile.uuid}</Card.Subtitle>
          <Card.Text>Email {profile.firstName} at <a href={`mailto:${profile.email}`}>{profile.email}</a></Card.Text>
        </Card.Body>
      </Card>
    );
  } else {
    if (thisIsUs) {
      profileHtml = (<p>You have not created a profile yet.</p>);
    } else {
      profileHtml = (<p>This user has not created a profile yet.</p>);
    }
  }

  return (
    <>
      {thisIsUs ? (<>
        <h3>Update Your Profile</h3>
        <Form onSubmit={ send }>
          <Form.Group controlId="formContent">
            <InputGroup className="mb-3">
              <InputGroup.Prepend>
                <InputGroup.Text>First and last name</InputGroup.Text>
              </InputGroup.Prepend>
              <FormControl name="firstName" value={profile?.firstName} placeholder={profile?.firstName ?? "Oski"} />
              <FormControl name="lastName" value={profile?.lastName} placeholder={profile?.firstName ?? "Bear"} />
            </InputGroup>

            <InputGroup className="mb-3">
              <FormControl
                value={profile?.email}
                placeholder={profile?.email ?? "oski@berkeley.edu"}
                name="email"
                type="email"
              />
            </InputGroup>
          </Form.Group>
          <Button variant="primary" type="submit">
            Update!
          </Button>
        </Form>

        <hr />

        <h3>Your Profile</h3>
        { profileHtml}
      </>) : (<>
        <h3>Profile</h3>
        { profileHtml }
      </>)}
    </>
  );
}

export default Profile;
