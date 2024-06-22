import {
    PHASE_DEVELOPMENT_SERVER,
    PHASE_PRODUCTION_BUILD,
} from 'next/constants';

/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    compiler: {
        removeConsole: process.env.NODE_ENV === 'production'
    }
};

const nextConfigFunction = async (phase) => {
    if (phase === PHASE_DEVELOPMENT_SERVER || phase === PHASE_PRODUCTION_BUILD) {
        const withPWA = (await import("@ducanh2912/next-pwa")).default({
            dest: "public",
            cacheOnFrontEndNav: true,
            reloadOnOnline: true,
            aggresiveFrontEndCache: true,
            disable: process.env.NODE_ENV === 'development',
            register: true,
            skipWaiting: true,
            swcMinify: true,
            workboxOptions: {
                disableDevLogs: true,
            }
        })

        return withPWA(nextConfig);
    }
    return nextConfig;
}

export default nextConfigFunction;
