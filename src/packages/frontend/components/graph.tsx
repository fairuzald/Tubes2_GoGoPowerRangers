"use client";

import {
  defineGraph,
  defineGraphConfig,
  defineLink,
  defineNode,
  GraphController,
  GraphLink,
  GraphNode,
} from "d3-graph-controller";
import "d3-graph-controller/default.css";
import React, { useEffect, useRef, useState } from "react";
import { GraphLinks, useQueryContext } from "./query-provider";

export type CustomType = "primary" | "secondary";

export interface CustomNode extends GraphNode<CustomType> {
  radius: number;
}

export interface CustomLink extends GraphLink<CustomType, CustomNode> {
  length: number;
}

const config = defineGraphConfig<CustomType, CustomNode, CustomLink>({
  nodeRadius: (node: CustomNode) => node.radius,
  simulation: {
    forces: {
      centering: {
        strength: (node: CustomNode) => (node.type === "primary" ? 0.5 : 0.1),
      },
      link: {
        length: (link: CustomLink) => link.length,
      },
    },
  },
});

const ForceGraph: React.FC = () => {
  const graphWrapperRef = useRef<HTMLDivElement>(null);

  const { state } = useQueryContext();

  let nodes: Record<string, CustomNode> = {};

  state.nodes.forEach((nodeId) => {
    const nodeGraph = defineNode<CustomType, CustomNode>({
      id: nodeId,
      type: "primary",
      isFocused: false,
      color: "#fca311",
      label: {
        color: "black",
        fontSize: "0.5rem",
        text: nodeId,
      },
      radius: 20,
    });

    nodes[nodeId] = nodeGraph;
  });

  let linkNodes: CustomLink[] = [];
  Object.values(state.linkNodes).forEach((link: GraphLinks) => {
    const { source, targets } = link;
    targets.forEach((dest) => {
      if (nodes[source] && nodes[dest]) {
        const link = defineLink<CustomType, CustomNode, CustomNode, CustomLink>(
          {
            source: nodes[source],
            target: nodes[dest],
            color: "black",
            label: {
              color: "black",
              fontSize: "1rem",
              text: "",
            },
            length: 200,
          }
        );

        linkNodes.push(link);
      }
    });
  });

  const graph = defineGraph<CustomType, CustomNode, CustomLink>({
    nodes: Object.values(nodes),
    links: linkNodes,
  });

  const [controller, setController] = useState<GraphController<
    CustomType,
    CustomNode,
    CustomLink
  > | null>(null);

  useEffect(() => {
    if (graphWrapperRef.current && graph) {
      const newController = new GraphController(
        graphWrapperRef.current,
        graph,
        config
      );
      setController(newController);
    }
  }, [graph]);

  return (
    <div
      ref={graphWrapperRef}
      style={{
        width: "100%",
        height: "250px",
        backgroundColor: "white",
      }}
    >
      {/* The graph will render inside this div */}
    </div>
  );
};

export default ForceGraph;
