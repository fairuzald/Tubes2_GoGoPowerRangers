"use client";
import { PathInfo } from "@/types/result";
import Image from "next/image";
import Link from "next/link";
import React from "react";
import { useQueryContext } from "./query-provider";

const Card: React.FC<{ data: PathInfo }> = ({ data }) => {
  return (
    <Link
      href={data.url}
      target="_blank"
      className={`container flex flex-row items-center rounded overflow-hidden shadow-lg bg-white hover:bg-gray-200 hover:cursor-pointer hover:shadow-2xl transition-shadow duration-200`}
    >
      <div>
        {/* Adjust the height as needed */}
        <Image
          src={data.image}
          alt={`Picture of ${data.title}`}
          className="rounded-l object-center object-cover"
          width={80}
          height={80}
        />
      </div>
      <div className="px-5">
        <div className="font-bold text-xl mb-2 text-black">{data.title}</div>
        <p className="text-gray-700 text-base">{data.description}</p>
      </div>
    </Link>
  );
};

const CardGrid: React.FC<{ data: PathInfo[] }> = ({ data }) => {
  return (
    <div className="p-5 bg-blue-200">
      <div className="grid grid-cols-1 gap-4">
        {data && data.length > 0 && data.map((card, index) => (
          <Card key={index} data={card} />
        ))}
      </div>
    </div>
  );
};

const CardGridResult   = () => {
  const { state } = useQueryContext();
  return (
    <section className="flex flex-wrap px-6 gap-6 w-full items-center justify-center">
      {
        state.result && state.result.length > 0 && state.result.map((path, index) => (
          <CardGrid key={index} data={path} />
        ))
      }
    </section>
  )
}

export default CardGridResult;
