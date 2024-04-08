/** @type {import('next').NextConfig} */
const nextConfig = {
    images: {
        remotePatterns: [
          {
            protocol: 'https',
            hostname: 'upload.wikimedia.org',
            port: '',
            pathname: '/**/*',
          },
        ],
      },
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: process.env.NEXT_API_URL+"/:path*",
      },
    ];
  },
};

export default nextConfig;
