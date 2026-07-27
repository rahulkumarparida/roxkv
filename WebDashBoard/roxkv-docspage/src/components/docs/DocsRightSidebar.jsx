import { useEffect, useState } from 'react';

export default function DocsRightSidebar() {
  const [headings, setHeadings] = useState([]);
  const [activeId, setActiveId] = useState('');

  useEffect(() => {
    const timer = setTimeout(() => {
      const elements = Array.from(document.querySelectorAll('main h2, main h3'));
      const headingData = elements.map((el) => {
        if (!el.id) {
          el.id = el.innerText.toLowerCase().replace(/[^a-z0-9]+/g, '-');
        }
        return {
          id: el.id,
          text: el.innerText,
          level: Number(el.tagName.charAt(1))
        };
      });
      setHeadings(headingData);
      
      // Initial active state from hash
      if (window.location.hash) {
        const id = window.location.hash.slice(1);
        setActiveId(id);
        const el = document.getElementById(id);
        if (el) {
          // Adjust scroll position for fixed header
          const y = el.getBoundingClientRect().top + window.scrollY - 80;
          window.scrollTo({ top: y, behavior: 'smooth' });
        }
      }
    }, 100);

    return () => clearTimeout(timer);
  }, []);

  useEffect(() => {
    if (headings.length === 0) return;

    const observer = new IntersectionObserver(
      (entries) => {
        // Find the topmost intersecting element
        const visibleEntries = entries.filter(e => e.isIntersecting);
        if (visibleEntries.length > 0) {
          setActiveId(visibleEntries[0].target.id);
        }
      },
      { rootMargin: '-80px 0% -60% 0%', threshold: 1.0 }
    );

    const elements = document.querySelectorAll('main h2, main h3');
    elements.forEach((el) => observer.observe(el));

    return () => observer.disconnect();
  }, [headings]);

  const handleClick = (e, id) => {
    e.preventDefault();
    setActiveId(id);
    
    // Update hash without jumping
    window.history.pushState(null, '', '#' + id);
    
    const el = document.getElementById(id);
    if (el) {
      // 80px offset to account for the sticky navbar
      const y = el.getBoundingClientRect().top + window.scrollY - 80;
      window.scrollTo({ top: y, behavior: 'smooth' });
    }
  };

  if (headings.length === 0) return null;

  return (
    <aside className="w-64 flex-shrink-0 hidden xl:block overflow-y-auto h-[calc(100vh-64px)] sticky top-16 pl-8 py-8">
      <h4 className="font-semibold text-[#f1f1f4] mb-4 text-xs uppercase tracking-widest opacity-60">
        On this page
      </h4>
      <ul className="space-y-2 text-sm">
        {headings.map((h) => (
          <li key={h.id} style={{ marginLeft: h.level === 3 ? '1rem' : '0' }}>
            <a
              href={`#${h.id}`}
              onClick={(e) => handleClick(e, h.id)}
              className={`block transition-colors duration-200 ${
                activeId === h.id ? 'text-[#a78bfa] font-medium' : 'text-[#9898ab] hover:text-[#f1f1f4]'
              }`}
            >
              {h.text}
            </a>
          </li>
        ))}
      </ul>
    </aside>
  );
}
