import map from "lodash/map";
import some from "lodash/some";
import range from "lodash/range";
import forEach from "lodash/forEach";
import debounce from "lodash/debounce";
import * as d3 from "d3";
import React, { Component } from "react";

import {
  GraphSvg,
  GraphWrapper,
  Instructions,
  Legend,
  LegendCircle,
  LegendItem,
  LegendLabel,
  ResetButton,
} from "./ui/results-graph-styles";

interface NodeDatum extends d3.SimulationNodeDatum {
  id: string;
  title: string;
  degree: number;
  x?: number;
  y?: number;
  fx?: number;
  fy?: number;
}

interface LinkDatum {
  source: string;
  target: string;
}

interface GraphProps {
  paths: Array<Array<{ title: string }>>;
}

interface GraphState {
  nodesData: NodeDatum[];
  linksData: LinkDatum[];
}

function getWikipediaPageUrl(title: string) {
  const baseUrl = "https://en.wikipedia.org/wiki/";
  const sanitizedPageTitle = title.replace(/ /g, "_");
  return `${baseUrl}${encodeURIComponent(sanitizedPageTitle)}`;
}

const DEFAULT_CHART_HEIGHT = 600;

class Graph extends Component<GraphProps, GraphState> {
  private graph: d3.Selection<SVGGElement, unknown, null, undefined> | null =
    null;
  private nodes: d3.Selection<
    SVGCircleElement,
    NodeDatum,
    SVGGElement,
    unknown
  > | null = null;
  private links: d3.Selection<
    SVGLineElement,
    LinkDatum,
    SVGGElement,
    unknown
  > | null = null;
  private zoomable: d3.Selection<
    SVGSVGElement,
    unknown,
    null,
    undefined
  > | null = null;
  private graphWidth: number | null = null;
  private nodeLabels: d3.Selection<
    SVGTextElement,
    NodeDatum,
    SVGGElement,
    unknown
  > | null = null;
  private simulation: d3.Simulation<NodeDatum, LinkDatum> | null = null;
  private ticksPerRender: number | null = null;
  // private graphRef: SVGSVGElement | null = null;
  private color: d3.ScaleOrdinal<string, string>;
  // private zoom: d3.ZoomBehavior<Element, unknown>;
  private zoom: d3.ZoomBehavior<SVGSVGElement, unknown>; 
  private debouncedResetGraph: () => void;
  private graphWrapperRef: React.RefObject<HTMLDivElement> = React.createRef();

  private graphRef: React.RefObject<SVGSVGElement>;

  constructor(props: GraphProps) {
    super(props);
    this.zoom = d3.zoom<SVGSVGElement, unknown>().on("zoom", (event) => this.zoomed(event));
  }

  getGraphWidth() {
    return this.graphWrapperRef.current?.getBoundingClientRect().width || 0;
  }

  /* Returns a list of nodes and a list of links which make up the graph. */
  getGraphData(): GraphState {
    const { paths } = this.props;

    const nodesData: NodeDatum[] = [];
    const linksData: LinkDatum[] = [];

    paths.forEach((path) => {
      path.forEach((node, i) => {
        const currentNodeId = node.title;

        // Add node if it has not yet been added by some other path.
        if (!some(nodesData, ["id", currentNodeId])) {
          nodesData.push({
            id: currentNodeId,
            title: node.title,
            degree: i,
          });
        }

        // Add link if this is not the start node.
        if (i !== 0) {
          linksData.push({
            source: path[i - 1].title,
            target: currentNodeId,
          });
        }
      });
    });

    return {
      nodesData,
      linksData,
    };
  }

  /* Returns a list of labels for use in the legend. */
  getLegendLabels() {
    const { paths } = this.props;
    const pathsLength = paths[0].length;

    return map(range(0, pathsLength), (i) => {
      if (i === 0 && pathsLength === 1) {
        return "Start / end page";
      } else if (i === 0) {
        return "Start page";
      } else if (i === pathsLength - 1) {
        return "End page";
      } else {
        const degreeOrDegrees = i === 1 ? "degree" : "degrees";
        return `${i} ${degreeOrDegrees} away`;
      }
    });
  }

  updateElementLocations(): void {
    if (
      !this.simulation ||
      !this.links ||
      !this.nodes ||
      !this.nodeLabels ||
      !this.ticksPerRender
    ) {
      // Exit if the simulation or DOM elements are not set up yet.
      return;
    }

    for (let i = 0; i < this.ticksPerRender; i++) {
      this.simulation.tick();
    }

    this.links
      .attr("x1", (d: any) => d.source.x)
      .attr("y1", (d: any) => d.source.y)
      .attr("x2", (d: any) => d.target.x)
      .attr("y2", (d: any) => d.target.y);

    this.nodes
      .attr("cx", (d: NodeDatum) => d.x ?? 0) // Use nullish coalescing to handle undefined x/y
      .attr("cy", (d: NodeDatum) => d.y ?? 0);

    this.nodeLabels.attr(
      "transform",
      (d: NodeDatum) => `translate(${d.x ?? 0}, ${d.y ?? 0})`
    );

    if (this.simulation.alpha() > 0.001) {
      requestAnimationFrame(() => this.updateElementLocations());
    }
  }

  /* Updates the zoom level of the graph when a zoom event occurs. */
  zoomed(event: d3.D3ZoomEvent<SVGGElement, unknown>): void {
    if (!this.graph) {
      return;
    }

    const transform = event.transform;
    this.graph.attr(
      "transform",
      `translate(${transform.x}, ${transform.y}) scale(${transform.k})`
    );
  }

  /* Drag started event. */
  dragstarted(
    event: d3.D3DragEvent<SVGCircleElement, NodeDatum, unknown>,
    d: NodeDatum
  ): void {
    if (!event.active) {
      this.simulation?.alphaTarget(0.3).restart();
      requestAnimationFrame(() => this.updateElementLocations());
    }
    d.fx = d.x;
    d.fy = d.y;
  }

  /* Dragged event. */
  dragged(
    event: d3.D3DragEvent<SVGCircleElement, NodeDatum, unknown>,
    d: NodeDatum
  ): void {
    d.fx = event.x;
    d.fy = event.y;
  }

  /* Drag ended event. */
  dragended(
    event: d3.D3DragEvent<SVGCircleElement, NodeDatum, unknown>,
    d: NodeDatum
  ): void {
    if (!event.active) {
      this.simulation?.alphaTarget(0).restart();
    }
    d.fx = undefined;
    d.fy = undefined;
  }

  componentDidMount() {
    const { paths } = this.props;
    const pathsLength = paths[0].length;
    const targetPageTitle = paths[0][pathsLength - 1].title;

    const { nodesData, linksData } = this.getGraphData();

    // Update the nubmer of ticks of the force simulation to run for each render according to how
    // many nodes will be rendered.
    this.ticksPerRender = 3 + Math.floor(nodesData.length / 20);

    this.graphWidth = this.getGraphWidth();

    this.zoomable = d3
      .select(this.graphRef!)
      .attr("width", "100%")
      .attr("height", "100%")
      .call(this.zoom);

    this.graph = this.zoomable.append("g");

    // Direction arrows.
    const defs = this.graph.append("defs");

    const markers = {
      arrow: 18,
      "arrow-end": 22,
    };

    forEach(markers, (refX, id) => {
      defs
        .append("marker")
        .attr("id", id)
        .attr("viewBox", "0 -5 10 10")
        .attr("refX", refX)
        .attr("refY", 0)
        .attr("markerWidth", 8)
        .attr("markerHeight", 8)
        .attr("orient", "auto")
        .append("svg:path")
        .attr("d", "M0,-5L10,0L0,5");
    });

    // Links.
    this.links = this.graph
      .append("g")
      .attr("class", "links")
      .selectAll("line")
      .data(linksData)
      .enter()
      .append("line")
      .attr("fill", "none")
      .attr("marker-end", (d) => {
        // Use a different arrow marker for links to the target page since it has a larger radius.
        if (d.target === targetPageTitle) {
          return "url(#arrow-end)";
        } else {
          return "url(#arrow)";
        }
      });

    // Nodes.
    this.nodes = this.graph
      .append("g")
      .attr("class", "nodes")
      .selectAll("circle")
      .data(nodesData)
      .enter()
      .append("circle")
      .attr("r", (d) => {
        if (d.degree === 0 || d.degree === pathsLength - 1) {
          return 10;
        } else {
          return 6;
        }
      })
      .attr("fill", (d) => this.color(d.degree))
      .attr("stroke", (d) => d3.rgb(this.color(d.degree)).darker(2))
      .call(
        d3
          .drag()
          .on("start", this.dragstarted.bind(this))
          .on("drag", this.dragged.bind(this))
          .on("end", this.dragended.bind(this))
      );

    // Node labels.
    this.nodeLabels = this.graph
      .append("g")
      .attr("class", "node-labels")
      .selectAll("text")
      .data(nodesData)
      .enter()
      .append("text")
      .attr("x", (d) => {
        if (d.degree === 0 || d.degree === pathsLength - 1) {
          return 14;
        } else {
          return 10;
        }
      })
      .attr("y", 4)
      .text((d) => d.title);

    // Open Wikipedia page when node is double clicked.
    this.nodes.on("click", (d) => {
      window.open(getWikipediaPageUrl(d.id), "_blank");
    });

    // Force simulation.
    this.simulation = d3
      .forceSimulation()
      .force(
        "link",
        d3.forceLink().id((d) => d.id)
      )
      .force("charge", d3.forceManyBody().strength(-300).distanceMax(500))
      .force(
        "center",
        d3.forceCenter(this.graphWidth / 2, DEFAULT_CHART_HEIGHT / 2)
      );

    this.simulation.nodes(nodesData);
    this.simulation.force("link").links(linksData);

    requestAnimationFrame(() => this.updateElementLocations());

    // Reset the graph on page resize.
    window.addEventListener("resize", this.debouncedResetGraph);
  }

  componentWillUnmount() {
    window.removeEventListener("resize", this.debouncedResetGraph);
  }

  /* Resets the graph to its original location. */
  resetGraph(forceReset) {
    const priorGraphWidth = this.graphWidth;
    this.graphWidth = this.getGraphWidth();

    if (forceReset || priorGraphWidth !== this.graphWidth) {
      // Update the center of the simulation force and restart it.
      this.simulation.force(
        "center",
        d3.forceCenter(this.graphWidth / 2, DEFAULT_CHART_HEIGHT / 2)
      );
      this.simulation.alpha(0.3).restart();

      // Update the width of the SVG and reset it.
      this.zoomable.attr("width", this.graphWidth);
      this.zoomable
        .transition()
        .duration(750)
        .call(this.zoom.transform, d3.zoomIdentity);

      requestAnimationFrame(() => this.updateElementLocations());
    }
  }

  /* Renders the legend. */
  renderLegend() {
    const legendContent = this.getLegendLabels().map((label, i) => {
      return (
        <LegendItem key={i}>
          <LegendCircle
            fill={this.color(i)}
            stroke={d3.rgb(this.color(i)).darker(2)}
          />
          <LegendLabel>{label}</LegendLabel>
        </LegendItem>
      );
    });
    return <Legend>{legendContent}</Legend>;
  }

  render() {
    return (
      <GraphWrapper className="graph-wrapper">
        {this.renderLegend()}

        <Instructions>
          <p>Drag to pan. Scroll to zoom.</p>
          <p>Click node to open Wikipedia page.</p>
        </Instructions>

        <ResetButton onClick={this.resetGraph.bind(this, true)}>
          <svg viewBox="0 0 100 100">
            <path d="m49.528 87h-0.06c-18.563 0-34.132-13.316-37.017-31.588-0.172-1.091-1.197-1.839-2.288-1.667s-1.836 1.201-1.664 2.292c3.195 20.225 20.422 34.963 40.963 34.963h0.066c11.585 0 22.714-4.672 30.542-12.814 7.451-7.751 11.311-17.963 10.869-28.751-0.952-23.211-19.169-41.394-41.474-41.394-15.237 0-29.288 8.546-36.465 21.722v-18.763c0-1.104-0.896-2-2-2s-2 0.896-2 2v25c0 1.104 0.896 2 2 2h25c1.104 0 2-0.896 2-2s-0.896-2-2-2h-20.635c6.034-13.216 19.456-21.961 34.101-21.961 20.152 0 36.613 16.497 37.476 37.557 0.397 9.688-3.067 18.861-9.755 25.818-7.078 7.361-17.156 11.586-27.659 11.586z" />
          </svg>
        </ResetButton>

        <GraphSvg ref={(r) => (this.graphRef = r)} />
      </GraphWrapper>
    );
  }
}

// TODO: add prop types
Graph.propTypes = {};

export default Graph;
