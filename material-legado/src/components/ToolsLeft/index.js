import React, { Component } from 'react';
import { Container } from './styles';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import PageTab from '../PageTab';
import ComponentTab from '../ComponentTab';
import StrutureTab from '../StructureTab';
import './tabs.css';

export default class ToolsLeft extends Component {
    render() {
        return (
            <div style={Container} className={this.props.className}>
                <Tabs>
                    <TabList>
                        <Tab>Páginas</Tab>
                        <Tab>Componentes</Tab>
                        <Tab>Estrutura</Tab>
                    </TabList>

                    <TabPanel>
                        <PageTab />
                    </TabPanel>
                    <TabPanel>
                        <ComponentTab />
                    </TabPanel>
                    <TabPanel>
                        <StrutureTab />
                    </TabPanel>
                </Tabs>
            </div>
        );
    }
}