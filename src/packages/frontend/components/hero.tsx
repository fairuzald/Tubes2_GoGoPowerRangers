// pages/index.tsx
import React from "react";
import Image from "next/image";
import InputQuery from "@/components/input-query";
import { Button } from "./ui/button";
// import { QueryProvider } from '@/context/QueryProvider'; // Make sure to create this context provider

export default function Home() {
  // Additional logic for handling form submission, if necessary
  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    // Implement what should happen on form submission
  };

  return (
    <div className="container mx-auto">
      <Image
        src="/wikirace.png"
        alt="logo wikirace"
        width={600}
        height={200}
        className="mx-auto"
      />

      <p className="text-center text-3xl font-bold mt-3">
        Find the shortest paths from
      </p>
    </div>
  );
}
