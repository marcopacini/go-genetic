import React from 'react';
import {
    Button,
    ButtonGroup,
    ButtonToolbar,
    Card,
    Col,
    Container,
    Dropdown,
    Image,
    Nav,
    Navbar,
    Row,
    Spinner
} from 'react-bootstrap';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

class Preview extends React.Component {
    constructor() {
        super();
        this.state = {
            started: false,
            src: "http://localhost:3000/placeholder-md.png",
            counter: 0,
            timer: "1000",
        };

        this.start = this.start.bind(this);
        this.update = this.update.bind(this);
        this.setTimer = this.setTimer.bind(this);
    }

    componentDidMount() {
        this.timer = setInterval(
            () => this.update(),
            parseInt(this.state.timer, 10)
        );
    }

    componentWillUnmount() {
        clearInterval(this.timer)
    }

    start() {
        fetch("http://localhost:3001/start")
            .then((response) => {
                this.setState(state => ({
                    started: !state.started
                }));
            })
    }

    update() {
        if (this.state.started) {
            this.setState((prevState) => {
                return {
                    src: "http://localhost:3001/best?n=" + this.state.counter,
                    counter: prevState.counter + 1
                }
            });
        }
    }

    setTimer(eventKey) {
        this.setState({
            timer: eventKey
        });

        clearInterval(this.timer);

        if (eventKey !== "0") {
            this.timer = setInterval(
                () => this.update(),
                parseInt(eventKey, 10)
            );
        }
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
                    </Card.Text>
                    <Navbar bg="light" expand="md">
                        <Navbar.Toggle aria-controls="basic-navbar-nav" />
                        <Nav className="mr-auto">
                            <ButtonToolbar>
                                <Button
                                    onClick={this.start}
                                    disabled={this.state.started}
                                    variant={!this.state.started ? "outline-success" : "outline-secondary"}
                                >
                                    {this.state.started ? <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true"/> : "Start"}
                                </Button>
                            </ButtonToolbar>
                        </Nav>
                        <Nav>
                            <Dropdown as={ButtonGroup}>
                                <Button
                                    variant="outline-secondary"
                                    onClick={this.update}
                                >
                                    Update Now
                                </Button>

                                <Dropdown.Toggle split variant="outline-secondary" id="dropdown-split-basic" />

                                <Dropdown.Menu>
                                    <Dropdown.Header>Auto Refresh</Dropdown.Header>
                                    <Dropdown.Item onSelect={this.setTimer} eventKey="1000" active={this.state.timer === "1000"}>1s</Dropdown.Item>
                                    <Dropdown.Item onSelect={this.setTimer} eventKey="5000" active={this.state.timer === "5000"}>5s</Dropdown.Item>
                                    <Dropdown.Item onSelect={this.setTimer} eventKey="30000" active={this.state.timer === "30000"}>30s</Dropdown.Item>
                                    <Dropdown.Divider />
                                    <Dropdown.Item onSelect={this.setTimer} eventKey="0" active={this.state.timer === "0"}>Off</Dropdown.Item>
                                </Dropdown.Menu>
                            </Dropdown>
                        </Nav>
                    </Navbar>

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
