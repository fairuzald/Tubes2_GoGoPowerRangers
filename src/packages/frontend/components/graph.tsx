"use client";
import React, { useEffect, useRef, useState, useMemo } from "react";
import {
  defineGraph,
  defineGraphConfig,
  defineLink,
  defineNode,
  GraphController,
  GraphNode,
  GraphLink,
} from "d3-graph-controller";
import "d3-graph-controller/default.css";
import { useQueryContext, GraphLinks } from "./query-provider";

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

  const { nodes, links } = useMemo(() => {
    let nodes: Record<string, CustomNode> = {};
    let links: CustomLink[] = [];

    state.nodes.forEach((nodeId) => {
      nodes[nodeId] = defineNode<CustomType, CustomNode>({
        id: nodeId,
        type: "primary",
        isFocused: false,
        color: "#fca311",
        label: {
          color: "black",
          fontSize: "0.5rem",
          text: nodeId,
        },
        radius: 32,
      });
    });

    Object.values(state.linkNodes).forEach((link: GraphLinks) => {
      link.targets.forEach((dest) => {
        if (nodes[link.source] && nodes[dest]) {
          links.push(
            defineLink<CustomType, CustomNode, CustomNode, CustomLink>({
              source: nodes[link.source],
              target: nodes[dest],
              color: "black",
              label: {
                color: "black",
                fontSize: "1rem",
                text: "",
              },
              length: 128,
            })
          );
        }
      });
    });

    return { nodes: Object.values(nodes), links };
  }, [state]);

  const [controller, setController] = useState<GraphController<
    CustomType,
    CustomNode,
    CustomLink
  > | null>(null);

  useEffect(() => {
    if (graphWrapperRef.current) {
      const graph = defineGraph<CustomType, CustomNode, CustomLink>({
        nodes,
        links,
      });
      const newController = new GraphController(
        graphWrapperRef.current,
        graph,
        config
      );
      setController(newController);
    }
  }, [nodes, links]);

  return (
    <div
      ref={graphWrapperRef}
      style={{
        width: "600px",
        height: "400px",
        backgroundColor: "white",
      }}
    ></div>
  );
};

export default ForceGraph;
