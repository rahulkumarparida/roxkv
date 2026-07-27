import React from 'react';
import { footer } from '../data/content';
import logo from '../assets/roxkvLogo.png';
import { Heart } from 'lucide-react';

export default function Footer() {
  return (
    <footer className="bg-[#070709] border-t border-[#1e1e2e]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
          <div>
            <div className="flex items-center gap-2">
              <img src={logo} alt="RoxAI Logo" className="h-7 w-auto" />
              <span className="text-xl font-bold text-white">RoxAI</span>
            </div>
            <p className="text-sm text-[#9898ab] mt-3">{footer.tagline}</p>
            <div className="text-xs text-[#5e5e73] mt-4 flex items-center gap-1">
              <Heart className="w-3 h-3 text-purple-400" /> {footer.madeBy}
            </div>
          </div>
          
          <div>
            <h4 className="text-sm font-semibold text-[#5e5e73] uppercase tracking-wider mb-4">Product</h4>
            {footer.links.product.map((link, index) => (
              <a 
                key={index} 
                href={link.href} 
                className="block mb-3 text-sm text-[#5e5e73] hover:text-purple-400 transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                {link.label}
              </a>
            ))}
          </div>
          
          <div>
            <h4 className="text-sm font-semibold text-[#5e5e73] uppercase tracking-wider mb-4">Connect</h4>
            {footer.links.connect.map((link, index) => (
              <a 
                key={index} 
                href={link.href} 
                className="block mb-3 text-sm text-[#5e5e73] hover:text-purple-400 transition-colors"
                target={link.href.startsWith('mailto:') ? "_self" : "_blank"}
                rel={link.href.startsWith('mailto:') ? "" : "noopener noreferrer"}
              >
                {link.label}
              </a>
            ))}
          </div>
        </div>
        
        <div className="mt-12 pt-8 border-t border-[#1e1e2e] flex flex-col sm:flex-row justify-between items-center gap-4">
          <p className="text-xs text-[#5e5e73]">{footer.copyright}</p>
          <p className="text-xs text-[#5e5e73] flex items-center gap-1">
            <Heart className="w-3 h-3 text-purple-400" /> {footer.madeWith}
          </p>
        </div>
      </div>
    </footer>
  );
}
