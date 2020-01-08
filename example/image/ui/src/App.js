import React from 'react';
import {Button, Card, Col, Container, Image, Row } from 'react-bootstrap';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

class Preview extends React.Component {
    constructor() {
        super();
        this.state = {
            started: false,
            src: "http://localhost:3001/best"
        };

        this.handleStart= this.handleStart.bind(this);
    }

    componentDidMount() {
        this.timerID = setInterval(
            () => this.tick(),
            1000
        );
    }

    componentWillUnmount() {
        clearInterval(this.timerID)
    }

    tick() {
        if (this.state.started) {
            this.setState({
                src: "http://localhost:3001/best?date=" + new Date()
            })
        }
    }

    handleStart() {
        fetch("http://localhost:3001/start")
            .then((response) => {
                this.setState(state => ({
                    started: !state.started
                }));
            })
    }

    render() {
        return (
            <Card>
                <Card.Header>Genetic Art, powered by <i>go-genetic</i></Card.Header>
                <Card.Body>
                    <Card.Text>
                        <b>n.b.</b> <i>working in progress</i>
                        <br /><br />
                        <Image src={this.state.src} rounded />
                        <br /><br />
                        <Button
                            onClick={this.handleStart}
                            disabled={this.state.started}
                            variant="primary"
                            size="xl"
                        >
                            Start
                        </Button>
                    </Card.Text>
                </Card.Body>
            </Card>
        );
    }
}

function App() {
  return (
    <Container className="h-100">
      <Row className="h-100 justify-content-center align-items-center text-center">
          <Col lg={6}>
            <Preview />
          </Col>
      </Row>
    </Container>
  );
}

export default App;
