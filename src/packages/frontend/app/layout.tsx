"use client";
import Navbar from "@/components/navbar";
import { QueryProvider } from "@/components/query-provider";
import { Poppins } from "next/font/google";
import React from "react";
import { Toaster } from "react-hot-toast";
import "./globals.css";

const poppins = Poppins({
  subsets: ["latin"],
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  display: "swap",
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [expandNavbar, setExpandNavbar] = React.useState(false);

  return (
    <html lang="en">
      <body className={poppins.className}>
        <Toaster />
        <QueryProvider>
          <Navbar
            expandNavbar={expandNavbar}
            setExpandNavbar={setExpandNavbar}
          />
          {children}
        </QueryProvider>
      </body>
    </html>
  );
}
