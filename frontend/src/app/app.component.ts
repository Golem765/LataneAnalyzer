import { Component, ElementRef, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import * as cytoscape from 'cytoscape';
import * as vis from 'vis';
import { User } from '../model/User';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  @ViewChild('graph') graphHolder: ElementRef;

  request: {
    beta: number,
    alpha: number,
    simulate: boolean,
    users: number
  } = {
    beta: 2,
    alpha: 2,
    simulate: true,
    users: 50,
  };

  response: {
    users: User[]
  };

  constructor(private http: HttpClient) {

  }

  analyze() {
    this.http.post('http://localhost:8000/go', this.request)
      .subscribe((data: any) => {
        this.response = data;
        console.log(data);
        //this.buildGraph();
        this.buildCytospaceGraph();
      });
  }

  buildGraph() {
    const nodes = new vis.DataSet([]);
    const edges = new vis.DataSet([]);
    let i = Math.random() * 500;
    let j = Math.random() * 500;
    for (const user of this.response.users) {
      nodes.add({
        x: i,
        y: j,
        id: user.id,
        title: user.name ? user.name : user.type,
        value: -user.pressure,
        shape: 'dot',
      });
      for (const mention in user.mentions) {
        edges.add({
          color: {
            color: 'red',
          },
          from: user.id,
          to: mention,
          arrows: 'from',
        });
      }
      if (this.request.simulate && user.friends) {
        for (const friend of user.friends) {
          edges.add({
            color: {
              color: 'blue',
            },
            from: user.id,
            to: friend,
          });
        }
      }
      i = i + 500 * Math.random();
      if (i > 2500) {
        i = Math.random() * 500;
        j = j + 500 * Math.random();
      }
    }
    const container = this.graphHolder.nativeElement;
    const data = {
      nodes: nodes,
      edges: edges,
    };
    const options = {
      physics: {
        stabilization: false,
        enabled: false,
      },
      nodes: {
        shape: 'dot',
        scaling: {
          customScalingFunction: function (min, max, total, value) {
            console.log(max);
            return value / max;
          },
          min: 5,
          max: 150,
        },
      },
    };
    new vis.Network(container, data, options);
  }

  buildCytospaceGraph() {
    const elements: any[] = [];
    let maxPressure = 1;
    let maxInfluence = 1;
    for (const user of this.response.users) {
      if (user.pressure === 0 && user.influence === 0) {
        continue;
      }
      elements.push({
        data: {
          id: user.id,
          weight: -user.pressure,
          name: user.name,
          link: user.link,
          influence: user.influence,
        },
      });
      maxPressure = Math.max(maxPressure, -user.pressure);
      maxInfluence = Math.max(maxInfluence, user.influence);

      for (const mention in user.mentions) {
        elements.push({
          data: {
            source: user.id, target: mention,
          },
          style: {
            'line-color': 'red',
            'width': 1,
            'curve-style': 'bezier',
            'source-arrow-shape': 'triangle',
            'source-arrow-color': 'red',
          },
        });
      }
      /*if (this.request.simulate && user.friends) {
        for (const friend of user.friends) {
          elements.push({
            data: {
              source: user.id, target: friend,
            },
          });
        }
      }*/
    }
    const cy = cytoscape({
      layout: {
        name: 'concentric',
        concentric: function (node: any) {
          return node.data('weight');
        },
      },
      style: [ // the stylesheet for the graph
        {
          selector: 'node',
          style: {
            'background-color': function (element) {
              const weight = element.data('weight');
              const influence = element.data('influence');
              if (influence > weight) {
                return 'blue';
              } else {
                return 'red';
              }
            },
            'background-opacity': function (element) {
              const weight = element.data('weight');
              const influence = element.data('influence');
              if (influence > weight) {
                return influence / maxInfluence;
              } else {
                return weight / maxPressure;
              }
            },
            'width': 100,
            'height': 40,
            'label': function (element) {
              if (element.data('name')) {
                return element.data('name');
              } else {
                return element.data('link');
              }
            },
            'text-halign': 'center',
            'text-valign': 'center',
            'font-size': '6px',
            'color': 'black',
            'shape': 'rectangle',
          },
        },
        {
          selector: 'edge',
          style: {
            'width': 1,
            'line-color': 'gray',
            'source-arrow-shape': 'triangle',
            'source-arrow-color': 'red',
          },
        },
      ],
      container: this.graphHolder.nativeElement,
      elements: elements,
    });
    cy.on('tap', 'node', function (evt) {
      const node = evt.target;
      alert(`pressure: ${node.data('weight')}, influence: ${node.data('influence')}`);
    });
  }
}
