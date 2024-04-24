"use client";
import {
  defineGraph,
  defineGraphConfig,
  defineLink,
  defineNode,
  GraphController,
  GraphLink,
  GraphNode
} from "d3-graph-controller";
import "d3-graph-controller/default.css";
import Link from "next/link";
import React, { useEffect, useRef, useState } from "react";
import { GraphLinks, useQueryContext } from "./query-provider";

export type CustomType = "primary" | "secondary";
const colors = ["#eb9834", "#3489eb", "#5634eb", "#5634eb", "#34eb46", "#4334eb"]

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
  const [controller, setController] = useState<GraphController<
    CustomType,
    CustomNode,
    CustomLink
  > | null>(null);

  useEffect(() => {
    if (!graphWrapperRef.current || !state) return;

    const nodes: Record<string, CustomNode> = {};
    Object.keys(state.nodes).forEach((nodeId: string) => {
      const depth = state.nodes[nodeId];
      const nodeGraph = defineNode<CustomType, CustomNode>({
        id: nodeId,
        type: "primary",
        isFocused: false,
        color:
          nodeId === state.selectedSource
            ? "#eb3a23"
            : nodeId === state.selectedDestination
              ? "#0c5925"
              : colors[depth],
        label: {
          color: "black",
          fontSize: "1rem",
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
          const link = defineLink<CustomType, CustomNode, CustomNode, CustomLink>({
            source: nodes[source],
            target: nodes[dest],
            color: "black",
            label: {
              color: "#FFFFFF",
              fontSize: "0.5rem",
              text: "",
            },
            length: 500,
          });

          linkNodes.push(link);
        }
      });
    });

    const graph = defineGraph<CustomType, CustomNode, CustomLink>({
      nodes: Object.values(nodes),
      links: linkNodes,
    });

    const newController = new GraphController(
      graphWrapperRef.current,
      graph,
      config
    );
    setController(newController);

    // Clean up controller on component unmount
    return () => {
      if (newController) {
        newController.shutdown(); // Stop the simulation before disposing
      }
    };
  }, [state]);

  return (
    Boolean(state.selectedSource) && Boolean(state.selectedDestination) && Boolean(state.nodes) && Boolean(state.result.length > 0) && (
      <>
        <div >
          {state.runtime !== undefined && state.articleCount !== undefined &&
            <p className="text-white text-center text-lg: lg:text-xl 2xl:text-2xl max-w-[1000px]">
              Found {" "}
              <span className="text-yellow-hover font-bold">{state.result.length} {state.result.length > 1 ? "paths" : "path"}</span>
              {" "}with{" "}
              <span className="text-yellow-hover font-bold">{state.result[0].length - 1} {state.result[0].length > 1 ? "degrees" : "degree"} of separation</span>
              {" "}  from{" "}
              <Link href={state.selectedSource} className="text-yellow-hover font-bold hover:underline hover:underline-offset-4">{state.source}</Link>
              {" "} to{" "}
              <Link href={state.selectedDestination} className="text-yellow-hover font-bold hover:underline hover:underline-offset-4">{state.destination}</Link>
              {" "} in{" "}
              <span className="text-yellow-hover font-bold">{state.runtime} seconds</span> runtime with <span className="text-yellow-hover font-bold">{state.articleCount} articles checked</span>
              {" "} and through {" "}
              <span className="text-yellow-hover font-bold">{(state.result[0].length-1)*state.result.length} articles</span>

            </p>
          }

        </div>
        <div
          ref={graphWrapperRef}
          style={{ width: "80%", height: "650px", }}
          className="bg-white rounded-2xl border-4 border-yellow-primary"
        />
      </>
    )
  )
};

export default ForceGraph;
