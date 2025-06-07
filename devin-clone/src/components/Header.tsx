import React from 'react';
import Link from 'next/link';

const Header = () => {
  return (
    <header className="bg-white shadow-md sticky top-0 z-50">
      <nav className="container mx-auto px-6 py-3 flex justify-between items-center">
        {/* Logo Placeholder */}
        <Link href="/" className="text-2xl font-bold text-gray-800">
          Logo
        </Link>

        {/* Primary Navigation */}
        <div className="hidden md:flex space-x-4 items-center">
          <Link href="/" className="text-gray-600 hover:text-gray-900">Home</Link>
          <Link href="/enterprise" className="text-gray-600 hover:text-gray-900">Enterprise</Link>
          <Link href="/pricing" className="text-gray-600 hover:text-gray-900">Pricing</Link>
          <Link href="/customers" className="text-gray-600 hover:text-gray-900">Customers</Link>

          {/* View More Dropdown */}
          <div className="relative group">
            <button className="text-gray-600 hover:text-gray-900 focus:outline-none">
              View more <span className="text-xs">▼</span>
            </button>
            <div className="absolute left-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 hidden group-hover:block">
              <Link href="/about" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">About us</Link>
              <Link href="/careers" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Careers</Link>
              <Link href="/blog" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Blog</Link>
              <Link href="/contact" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Contact us</Link>
              <Link href="/docs" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Docs</Link>
              <Link href="/deepwiki" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">DeepWiki</Link>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex items-center space-x-3">
          <Link href="/login" className="text-gray-600 hover:text-gray-900">Login</Link>
          <button className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md">
            Get Started
          </button>
        </div>

        {/* Mobile Menu Button (placeholder) */}
        <div className="md:hidden">
          <button className="text-gray-600 hover:text-gray-900 focus:outline-none">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16m-7 6h7"></path>
            </svg>
          </button>
        </div>
      </nav>
    </header>
  );
};

export default Header;
