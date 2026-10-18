import React, { Component } from 'react';
import { Container } from './styles';
import DateInfo from '../DateInfo';
import LotericaInfo from '../LotericaInfo';
import InfraInfo from '../InfraInfo';
import { AppConsumer } from '../../service/AppContext';
import PageElementPreview from '../PageElementPreview';
import './page.css';

export default class Preview extends Component {

  render() {
    return (
      <AppConsumer>
        {data => {
          return (
            <div style={Container} className={this.props.className} >
              {data.selectedPage &&
                (
                  <React.Fragment>
                    <div id="page-preview" className="page">
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
                          <PageElementPreview
                            element={component}
                            data={data}
                            key={index} />
                        ))
                      }

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
