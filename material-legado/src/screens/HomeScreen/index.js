import React, { Component } from 'react';
import './styles.css';
import Preview from '../../components/Preview';
import Design from '../../components/Design';
import Header from '../../components/Header';
import ToolsLeft from '../../components/ToolsLeft';
import ToolsRight from '../../components/ToolsRight';
import Code from '../../components/Code';
import { AppConsumer } from '../../service/AppContext';

export default class HomeScreen extends Component {
  render() {
    return (
      <AppConsumer>
        {props => (
          <div className={props.viewMode === 'design' ? 'grid-container-design' : 'grid-container'}>
            <Header className="grid-header"></Header>
            {
              props.viewMode === 'design' && <ToolsLeft className="grid-tools"></ToolsLeft>
            }
            {props.viewMode === 'view' && <Preview className="grid-children"></Preview>}
            {props.viewMode === 'design' && <Design className="grid-design"></Design>}
            {props.viewMode === 'code' && <Code className="grid-children"></Code>}
            {
              props.viewMode === 'design' && <ToolsRight className="grid-configuration"></ToolsRight>
            }
          </div>
        )}
      </AppConsumer>
    );
  }
}
