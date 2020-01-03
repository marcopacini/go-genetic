import React from 'react';
import { Button, Card, Col, Container, Image, Row } from 'react-bootstrap';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  return (
    <Container className="h-100">
      <Row className="h-100 justify-content-center align-items-center text-center">
          <Col lg={6}>
          <Card>
            <Card.Header>Genetic Art, powered by <i>go-genetic</i></Card.Header>
            <Card.Body>
              <Card.Text>
                <b>n.b.</b> <i>working in progress</i>
                <br /><br />
                <Image src="http://localhost:3001" rounded />
                <br /><br />
              </Card.Text>
              <Button variant="primary">
                Generate
              </Button>
            </Card.Body>
          </Card>
          </Col>
      </Row>
    </Container>
  );
}

export default App;
