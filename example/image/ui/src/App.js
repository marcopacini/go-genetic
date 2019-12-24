import React from 'react';
import { Container, Row, Card } from 'react-bootstrap';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  return (
    <Container className="h-100">
      <Row className="h-100 justify-content-center align-items-center text-center">
          <Card className="w-50">
            <Card.Body>Genetic Art, powered by <i>go-genetic</i></Card.Body>
          </Card>
      </Row>
    </Container>
  );
}

export default App;
