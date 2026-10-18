import React, { Component } from 'react';
import { Container, InfoContainer, InfoItens } from './styles';
import DateInfo from '../DateInfo';
import LotericaInfo from '../LotericaInfo';
import InfraInfo from '../InfraInfo';
import { AppConsumer } from '../../service/AppContext';
import PageElement from '../PageElement';
import ListComponents from '../ListComponents';
import './page.css';

export default class Design extends Component {

  constructor(props) {
    super(props);
    this.state = {
      currentElement: null,
      cursor: {
        x: 0.0,
        y: 0.0
      }
    };
    this.pageRef = React.createRef();
  }

  keyDown = (event) => {
    const isPage = event.target.tagName.toUpperCase() === 'BODY';
    const keyCode = event.keyCode;
    if (keyCode === 37 || keyCode === 38 || keyCode === 39 || keyCode === 40 || keyCode === 46) {
      if(isPage){
        event.preventDefault();
      }
    }
    const selected = this.selected;
    if (selected && isPage) {
      if (keyCode === 46) {
        // Delete
        this.removeComponent(selected.id);
      }
      if (keyCode === 40) {
        // Down
        selected.properties.top++;
        this.updateValue('selected', selected);
      } else if (keyCode === 38) {
        // Up
        selected.properties.top--;
        this.updateValue('selected', selected);
      } else if (keyCode === 39) {
        // Right
        selected.properties.left++;
        this.updateValue('selected', selected);
      } else if (keyCode === 37) {
        // Left
        selected.properties.left--;
        this.updateValue('selected', selected);
      }
    }
  }

  componentDidMount = () => {
    document.addEventListener('keydown', this.keyDown);
  }

  componentWillUnmount = () => {
    document.removeEventListener('keydown', this.keyDown);
  }

  onDragOver = (event) => {
    event.preventDefault();
    this.updateXY(event);
  }

  onDrop = (event) => {
    if(event.preventDefault) { event.preventDefault(); }
    if(event.stopPropagation) { event.stopPropagation(); }
    const idComponent = event.dataTransfer.getData("text/plain");
    const data = ListComponents.find(item => item.id === idComponent);
    if (data) {
      const componente = new data.component({ x: this.state.cursor.x, y: this.state.cursor.y });
      this.addComponent(componente);
    }
    return false;
  }

  onMouseMove = (event) => {
    this.updateXY(event);
  }

  updateXY = (event) => {
    if (!this.pageRef || !this.pageRef.current) {
      return;
    }
    let rect = this.pageRef.current.getBoundingClientRect();
    let x = parseInt(event.clientX - rect.left, 10);
    let y = parseInt(event.clientY - rect.top, 10);
    if (x !== this.state.cursor.x || y !== this.state.cursor.y) {
      this.setState({ cursor: { x, y } }, this.updateElement);
    }
  }

  onMouseDown = (element) => {
    this.setState({ currentElement: element }, () => {
      this.updateValue('selected', element);
    });
  }

  onMouseUp = () => {
    this.setState({ currentElement: null });
  }

  updateElement = () => {
    if (!!this.state.currentElement) {
      this.setState(prevState => {
        const currentElement = prevState.currentElement;
        currentElement.properties.top = this.state.cursor.y;
        currentElement.properties.left = this.state.cursor.x;
        return { currentElement };
      }, () => {
        this.notifyUpdate();
      });
    }
  }

  render() {
    return (
      <AppConsumer>
        {data => {
          this.addComponent = data.addComponent;
          this.updateValue = data.updateValue;
          this.selected = data.selected;
          this.removeComponent = data.removeComponent;
          this.notifyUpdate = data.notifyUpdate;
          return (
            <div style={Container} className={this.props.className} >
              {data.selectedPage &&
                (
                  <React.Fragment>
                    <div ref={this.pageRef} id="page-preview" className="page" onMouseUp={this.onMouseUp} onDragOver={this.onDragOver} onDrop={this.onDrop} onMouseMove={this.onMouseMove}>
                      <div className="page-header">
                        <div className="page-header-date">
                          <DateInfo />
                        </div>
                        <div className="page-header-user">
                          <LotericaInfo />
                        </div>
                        <div className="page-header-info">
                          <InfraInfo />
                        </div>
                      </div>
                      {data.selectedPage && data.selectedPage.showTitle && <div className="title-container"><h1 className="title">{data.selectedPage.title}</h1></div>}
                      {
                        data.selectedPage && data.selectedPage.components && data.selectedPage.components.map((component, index) => (
                          <PageElement
                            onClick={(element) => { this.updateValue('selected', element) }}
                            element={component}
                            key={index}
                            removeComponent={data.removeComponent}
                            onMouseDown={this.onMouseDown}
                            onMouseUp={this.onMouseUp}
                            selected={data.selected} />
                        ))
                      }

                    </div>

                    <div style={InfoContainer}>
                      <div style={InfoItens}>
                        x: {this.state.cursor.x}
                        y: {this.state.cursor.y}
                      </div>
                    </div>
                  </React.Fragment>
                )
              }
            </div>
          )
        }}
      </AppConsumer>
    );
  }
}
