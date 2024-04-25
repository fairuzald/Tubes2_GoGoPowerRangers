"use client";
import { PathInfo } from "@/types/result";
import Image from "next/image";
import Link from "next/link";
import React from "react";
import { useQueryContext } from "./query-provider";

const Card: React.FC<{ data: PathInfo }> = ({ data }) => {
  return (
    <Link
      href={data.url || "/"}
      target="_blank"
      className={`container flex gap-5 px-3 py-3 flex-row items-center rounded overflow-hidden shadow-lg bg-white hover:bg-gray-200 hover:cursor-pointer hover:shadow-2xl transition-shadow duration-200`}
    >
      <div>
        {/* Adjust the height as needed */}
        <Image
          src={data.image || "/default.png"}
          alt={`Picture of ${data.title}`}
          width={80}
          height={80}
          className="rounded object-center object-cover w-[80px] h-[80px]"
        />
      </div>
      <div className="">
        <div className="font-bold text-xl mb-2 text-black">{data.title}</div>
        <p className="text-gray-700 text-base">{data.description}</p>
      </div>
    </Link>
  );
};

const CardGrid: React.FC<{ data: PathInfo[] }> = ({ data }) => {
  return (
    <div className="p-5 bg-blue-200 max-w-[500px] h-auto">
      <div className="flex flex-col gap-4 items-center justify-center h-full">
        {data &&
          data.length > 0 &&
          data.map((card, index) => <Card key={index} data={card} />)}
      </div>
    </div>
  );
};

const CardGridResult = () => {
  const { state } = useQueryContext();
  return (
    state.result &&
    state.result.length > 0 && (
      <section className="flex flex-col gap-10">
        <h3 className="text-lg lg:text-xl 2xl:text-2xl text-yellow-hover text-center font-bold">
          Individual paths
        </h3>
        <div className="flex flex-wrap px-6 gap-6 w-full items-stretch justify-center">
          {state.result.map((path, index) => (
            <CardGrid key={index} data={path} />
          ))}
        </div>
      </section>
    )
  );
};

export default CardGridResult;
