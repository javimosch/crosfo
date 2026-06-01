# DESIGN.md - Design System & UI Guidelines

## Design Philosophy

**CLI-First, Web-Enhanced, Premium Quality**

This design system reflects the core philosophy of the project: CLI-native applications with optional web interfaces. The visual language prioritizes functionality, clarity, and developer experience while maintaining aesthetic consistency across both terminal and web interfaces, elevated to premium quality standards.

### Core Principles

1. **Functional Minimalism**: Every visual element serves a purpose
2. **Developer-Centric**: Designed for technical users who value efficiency
3. **Terminal-Inspired**: Visual language references CLI aesthetics
4. **Performance-First**: Lightweight, fast-loading, no unnecessary bloat
5. **Consistency**: Unified design language across CLI and web interfaces
6. **Premium Quality**: Elevated aesthetics without sacrificing functionality

## Visual Identity

### Brand Attributes

- **Technical**: Precision, engineering, systems
- **Efficient**: Fast, streamlined, purposeful
- **Reliable**: Stable, trustworthy, dependable
- **Modern**: Contemporary, clean, forward-thinking
- **Premium**: Elevated aesthetics, refined details, professional polish

### Design Language

The visual language draws inspiration from:
- Terminal interfaces and command-line tools
- Developer tools and IDEs
- System monitoring dashboards
- Technical documentation

## Color System

### Primary Palette

```css
/* Dark Mode (Default) */
--color-bg-primary: #0d1117;      /* GitHub dark background */
--color-bg-secondary: #161b22;    /* Slightly lighter for cards */
--color-bg-tertiary: #21262d;     /* Borders, dividers */
--color-text-primary: #c9d1d9;    /* Primary text */
--color-text-secondary: #8b949e;  /* Secondary text */
--color-text-muted: #6e7681;      /* Muted text */
--color-accent-primary: #58a6ff;  /* Blue accent */
--color-accent-secondary: #238636; /* Green success */
--color-accent-danger: #da3633;   /* Red error */
--color-accent-warning: #d29922;  /* Yellow warning */
```

### Light Mode (Optional)

```css
--color-bg-primary: #ffffff;
--color-bg-secondary: #f6f8fa;
--color-bg-tertiary: #e1e4e8;
--color-text-primary: #24292f;
--color-text-secondary: #57606a;
--color-text-muted: #6e7681;
--color-accent-primary: #0969da;
--color-accent-secondary: #1a7f37;
--color-accent-danger: #cf222e;
--color-accent-warning: #9a6700;
```

### Semantic Colors

```css
/* Status Indicators */
--color-status-running: #238636;    /* Green */
--color-status-stopped: #da3633;     /* Red */
--color-status-paused: #d29922;     /* Yellow */
--color-status-unknown: #8b949e;    /* Gray */

/* Data Visualization */
--color-chart-1: #58a6ff;
--color-chart-2: #238636;
--color-chart-3: #a371f7;
--color-chart-4: #d29922;
--color-chart-5: #f85149;
```

## Typography

### Font Stack

```css
/* Primary Font - System UI */
--font-family-primary: -apple-system, BlinkMacSystemFont, "Segoe UI", 
                       "Noto Sans", Helvetica, Arial, sans-serif;

/* Monospace Font - Code/Technical */
--font-family-mono: "SF Mono", "Segoe UI Mono", "Roboto Mono", 
                    "Menlo", "Consolas", monospace;
```

### Type Scale

```css
/* Headings */
--font-size-h1: 2.5rem;      /* 40px */
--font-size-h2: 2rem;        /* 32px */
--font-size-h3: 1.5rem;      /* 24px */
--font-size-h4: 1.25rem;     /* 20px */

/* Body */
--font-size-body: 1rem;      /* 16px */
--font-size-small: 0.875rem; /* 14px */
--font-size-xs: 0.75rem;     /* 12px */

/* Monospace */
--font-size-mono: 0.875rem;  /* 14px */
```

### Font Weights

```css
--font-weight-normal: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;
```

### Line Heights

```css
--line-height-tight: 1.25;
--line-height-normal: 1.5;
--line-height-relaxed: 1.75;
```

## Spacing System

### Scale

```css
--spacing-0: 0;
--spacing-1: 0.25rem;  /* 4px */
--spacing-2: 0.5rem;   /* 8px */
--spacing-3: 0.75rem;  /* 12px */
--spacing-4: 1rem;     /* 16px */
--spacing-5: 1.25rem;  /* 20px */
--spacing-6: 1.5rem;   /* 24px */
--spacing-8: 2rem;     /* 32px */
--spacing-10: 2.5rem;  /* 40px */
--spacing-12: 3rem;    /* 48px */
--spacing-16: 4rem;    /* 64px */
```

### Usage Guidelines

- **--spacing-2**: Tight spacing, related elements
- **--spacing-4**: Default spacing, standard padding
- **--spacing-6**: Section spacing, comfortable breathing room
- **--spacing-8**: Major sections, distinct visual separation
- **--spacing-12**: Large sections, page-level spacing

## Component Design

### Buttons

#### Primary Button
```css
.btn-primary {
  background-color: var(--color-accent-primary);
  color: #ffffff;
  border: none;
  border-radius: 6px;
  padding: var(--spacing-2) var(--spacing-4);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.btn-primary:hover {
  background-color: #409eff;
}
```

#### Secondary Button
```css
.btn-secondary {
  background-color: var(--color-bg-tertiary);
  color: var(--color-text-primary);
  border: 1px solid var(--color-bg-tertiary);
  border-radius: 6px;
  padding: var(--spacing-2) var(--spacing-4);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.btn-secondary:hover {
  background-color: var(--color-bg-secondary);
}
```

#### Danger Button
```css
.btn-danger {
  background-color: var(--color-accent-danger);
  color: #ffffff;
  border: none;
  border-radius: 6px;
  padding: var(--spacing-2) var(--spacing-4);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.btn-danger:hover {
  background-color: #b62324;
}
```

### Cards

```css
.card {
  background-color: var(--color-bg-secondary);
  border: 1px solid var(--color-bg-tertiary);
  border-radius: 8px;
  padding: var(--spacing-6);
  margin-bottom: var(--spacing-6);
}

.card-header {
  font-size: var(--font-size-h4);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-4);
  color: var(--color-text-primary);
}

.card-body {
  color: var(--color-text-secondary);
  line-height: var(--line-height-normal);
}
```

### Status Indicators

```css
.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-running .status-dot {
  background-color: var(--color-status-running);
}

.status-stopped .status-dot {
  background-color: var(--color-status-stopped);
}
```

### Code Blocks

```css
.code-block {
  background-color: var(--color-bg-primary);
  border: 1px solid var(--color-bg-tertiary);
  border-radius: 6px;
  padding: var(--spacing-4);
  font-family: var(--font-family-mono);
  font-size: var(--font-size-mono);
  color: var(--color-text-primary);
  overflow-x: auto;
}

.inline-code {
  background-color: var(--color-bg-tertiary);
  border-radius: 4px;
  padding: 2px 6px;
  font-family: var(--font-family-mono);
  font-size: var(--font-size-small);
  color: var(--color-text-primary);
}
```

### Forms

```css
.form-group {
  margin-bottom: var(--spacing-6);
}

.form-label {
  display: block;
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
  margin-bottom: var(--spacing-2);
}

.form-input {
  width: 100%;
  padding: var(--spacing-3) var(--spacing-4);
  background-color: var(--color-bg-primary);
  border: 1px solid var(--color-bg-tertiary);
  border-radius: 6px;
  color: var(--color-text-primary);
  font-family: var(--font-family-primary);
  font-size: var(--font-size-body);
  transition: border-color 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}
```

## Layout Patterns

### Container

```css
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-6);
}

.container-narrow {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 var(--spacing-6);
}
```

### Grid System

```css
.grid {
  display: grid;
  gap: var(--spacing-6);
}

.grid-2 {
  grid-template-columns: repeat(2, 1fr);
}

.grid-3 {
  grid-template-columns: repeat(3, 1fr);
}

.grid-4 {
  grid-template-columns: repeat(4, 1fr);
}

@media (max-width: 768px) {
  .grid-2, .grid-3, .grid-4 {
    grid-template-columns: 1fr;
  }
}
```

### Flexbox Utilities

```css
.flex {
  display: flex;
}

.flex-col {
  flex-direction: column;
}

.items-center {
  align-items: center;
}

.justify-between {
  justify-content: space-between;
}

.gap-2 { gap: var(--spacing-2); }
.gap-4 { gap: var(--spacing-4); }
.gap-6 { gap: var(--spacing-6); }
```

## Animation & Transitions

### Timing Functions

```css
--transition-fast: 0.15s ease;
--transition-normal: 0.2s ease;
--transition-slow: 0.3s ease;
```

### Hover Effects

```css
.hover-lift {
  transition: transform var(--transition-normal) ease,
              box-shadow var(--transition-normal) ease;
}

.hover-lift:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
```

### Fade Animations

```css
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-in {
  animation: fadeIn var(--transition-slow) ease;
}
```

## Responsive Design

### Breakpoints

```css
--breakpoint-sm: 640px;
--breakpoint-md: 768px;
--breakpoint-lg: 1024px;
--breakpoint-xl: 1280px;
```

### Mobile-First Approach

- Default styles target mobile devices
- Use `min-width` media queries for larger screens
- Ensure touch targets are at least 44x44px
- Optimize for vertical scrolling on mobile

## Accessibility

### Color Contrast

- Ensure text contrast ratio of at least 4.5:1
- Use semantic colors for status indicators
- Provide text alternatives for color-coded information

### Keyboard Navigation

- All interactive elements must be keyboard accessible
- Provide visible focus indicators
- Support tab navigation in logical order
- Include skip-to-content links

### Screen Reader Support

- Use semantic HTML elements
- Provide ARIA labels for interactive elements
- Ensure form inputs have associated labels
- Include alt text for images

## Iconography

### Icon Style

- Use consistent stroke width (2px)
- Prefer outline icons over filled icons
- Maintain 24x24px base size
- Use SVG format for scalability

### Icon Sources

- [Lucide Icons](https://lucide.dev/) - Primary icon set
- [Heroicons](https://heroicons.com/) - Alternative option
- Custom SVG icons for specific needs

## Web UI Structure

### Page Layout

```
┌─────────────────────────────────────┐
│ Header (Logo, Navigation)           │
├─────────────────────────────────────┤
│                                     │
│ Main Content Area                   │
│ - Hero Section                      │
│ - Status Dashboard                  │
│ - Control Panel                    │
│ - Configuration                     │
│                                     │
├─────────────────────────────────────┤
│ Footer (Copyright, Links)          │
└─────────────────────────────────────┘
```

### Dashboard Components

1. **Status Overview**: Running/stopped status, uptime, resource usage
2. **Quick Actions**: Start/stop/restart buttons
3. **Configuration Panel**: Port, host, logging settings
4. **Log Viewer**: Real-time log display with filtering
5. **API Documentation**: Endpoint documentation and testing

### CLI-Web Consistency

- Use same terminology in CLI and web UI
- Maintain consistent color schemes
- Provide CLI command examples in web UI
- Include web UI URLs in CLI help text

## Premium Design Enhancements

### Typography Premium

**Extended Font Weight Scale**
```css
--font-weight-light: 300;      /* For subtle labels */
--font-weight-normal: 400;    /* Body text */
--font-weight-medium: 500;    /* Descriptions */
--font-weight-semibold: 600;  /* Subheaders */
--font-weight-bold: 700;      /* Headers */
```

**Letter-Spacing System**
```css
--letter-spacing-tight: -0.04em;  /* Large headers */
--letter-spacing-normal: -0.02em; /* Standard headers */
--letter-spacing-loose: 0.02em;   /* Small text */
```

**Orphan Prevention**
```css
text-wrap: balance;  /* Prevents single words on last line */
text-wrap: pretty;   /* Optimizes line breaks */
```

### Color Premium

**Enhanced Palette**
```css
/* Premium Backgrounds */
--color-bg-primary: #FAFAFA;      /* Off-white, not clinical */
--color-bg-secondary: #F5F5F5;    /* Slightly darker */
--color-bg-tertiary: #E5E5E5;     /* Borders */
--color-bg-elevated: #FFFFFF;     /* Cards, elevated elements */

/* Premium Text */
--color-text-primary: #1A1A1A;    /* Charcoal, not harsh black */
--color-text-secondary: #666666;  /* Muted gray */
--color-text-tertiary: #999999;   /* Subtle text */

/* Premium Accent */
--color-accent-primary: #00B4D8;  /* Teal accent */
--color-accent-hover: #0095B8;    /* Darker teal */
--color-accent-subtle: rgba(0, 180, 216, 0.1); /* Tinted shadows */
```

**Tinted Shadows**
```css
--shadow-light: rgba(0, 180, 216, 0.1);
--shadow-medium: rgba(0, 180, 216, 0.15);
--shadow-dark: rgba(0, 180, 216, 0.2);
```

### Layout Premium

**Variable Border-Radius**
```css
--radius-sm: 4px;   /* Small elements, inputs */
--radius-md: 8px;   /* Standard cards, buttons */
--radius-lg: 12px;  /* Large cards, modals */
--radius-xl: 16px;  /* Hero sections, major containers */
```

**Visual Depth System**
```css
--depth-1: 0 2px 8px var(--shadow-light);    /* Subtle elevation */
--depth-2: 0 4px 16px var(--shadow-medium); /* Standard elevation */
--depth-3: 0 8px 24px var(--shadow-dark);   /* Strong elevation */
```

**Asymmetric Spacing**
```css
--spacing-top-sm: 60px;
--spacing-bottom-sm: 80px;  /* More bottom padding */
--spacing-top-md: 80px;
--spacing-bottom-md: 100px;
```

### Interactivity Premium

**Hover States**
```css
.btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--depth-2);
}

.btn:active {
  transform: translateY(0) scale(0.98);
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: var(--depth-3);
}
```

**Focus States**
```css
:focus-visible {
  outline: 2px solid var(--color-accent-primary);
  outline-offset: 2px;
}

input:focus, button:focus {
  box-shadow: 0 0 0 3px var(--color-accent-subtle);
}
```

**Transitions**
```css
--transition-fast: 0.15s ease;
--transition-normal: 0.2s ease;
--transition-slow: 0.3s ease;
```

### Animation Premium

**Staggered Entry**
```css
@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stagger-1 { animation-delay: 0.1s; }
.stagger-2 { animation-delay: 0.2s; }
.stagger-3 { animation-delay: 0.3s; }
```

**Glassmorphism**
```css
.glass {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 
    0 4px 24px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
}
```

**Texture Overlays**
```css
.texture-noise {
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)' opacity='0.03'/%3E%3C/svg%3E");
  pointer-events: none;
}
```

### Component Premium

**Premium Buttons**
```css
.btn {
  border-radius: var(--radius-md);
  transition: all var(--transition-normal) ease;
  position: relative;
  overflow: hidden;
}

.btn::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent);
  transform: translateX(-100%);
  transition: transform var(--transition-slow) ease;
}

.btn:hover::after {
  transform: translateX(100%);
}
```

**Premium Cards**
```css
.card {
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-bg-tertiary);
  border-radius: var(--radius-lg);
  padding: var(--spacing-6);
  box-shadow: var(--depth-1);
  transition: all var(--transition-normal) ease;
  position: relative;
}

.card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--color-accent-primary), transparent);
  transform: translateX(-100%);
  transition: transform var(--transition-slow) ease;
}

.card:hover::before {
  transform: translateX(100%);
}
```

**Premium Forms**
```css
.form-input {
  border-radius: var(--radius-sm);
  transition: all var(--transition-normal) ease;
}

.form-input:focus {
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 3px var(--color-accent-subtle);
}

.form-input.error {
  border-color: var(--color-accent-danger);
  box-shadow: 0 0 0 3px rgba(218, 54, 51, 0.3);
}
```

### States Premium

**Loading States**
```css
.skeleton {
  background: linear-gradient(90deg, #E5E5E5 25%, #F5F5F5 50%, #E5E5E5 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: var(--radius-sm);
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
```

**Error States**
```css
.error-message {
  background: #FFF3F3;
  border: 1px solid #FFC7C7;
  color: #D32F2F;
  padding: var(--spacing-3) var(--spacing-4);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-small);
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}
```

**Empty States**
```css
.empty-state {
  text-align: center;
  padding: var(--spacing-12) var(--spacing-6);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-lg);
  border: 2px dashed var(--color-bg-tertiary);
}
```

### Accessibility Premium

**Skip-to-Content**
```css
.skip-to-content {
  position: absolute;
  top: -40px;
  left: 0;
  background: var(--color-accent-primary);
  color: white;
  padding: var(--spacing-2) var(--spacing-4);
  z-index: 10000;
  transition: top var(--transition-normal) ease;
}

.skip-to-content:focus {
  top: 0;
}
```

**Smooth Scroll**
```css
html {
  scroll-behavior: smooth;
  scroll-padding-top: 80px;
}
```

## Implementation Guidelines

### CSS Organization

```css
/* 1. CSS Variables (Design Tokens) */
:root {
  /* Colors, typography, spacing */
}

/* 2. Base Styles */
* {
  box-sizing: border-box;
}

body {
  font-family: var(--font-family-primary);
  color: var(--color-text-primary);
  background-color: var(--color-bg-primary);
}

/* 3. Layout Components */
.container, .grid, .flex { }

/* 4. UI Components */
.btn, .card, .form-input { }

/* 5. Utilities */
.text-center, .mt-4, .p-6 { }

/* 6. Responsive */
@media (max-width: 768px) { }
```

### JavaScript Integration

- Use vanilla JavaScript for interactivity
- Minimize external dependencies
- Implement graceful degradation
- Provide loading states for async operations

### Performance Optimization

- Minimize CSS file size
- Use system fonts to avoid loading delays
- Implement lazy loading for non-critical resources
- Optimize images and assets
- Use CSS animations instead of JavaScript when possible

## Design Tokens Export

### JSON Format

```json
{
  "colors": {
    "bg": {
      "primary": "#0d1117",
      "secondary": "#161b22",
      "tertiary": "#21262d"
    },
    "text": {
      "primary": "#c9d1d9",
      "secondary": "#8b949e",
      "muted": "#6e7681"
    },
    "accent": {
      "primary": "#58a6ff",
      "secondary": "#238636",
      "danger": "#da3633",
      "warning": "#d29922"
    }
  },
  "typography": {
    "fontFamily": {
      "primary": "-apple-system, BlinkMacSystemFont, sans-serif",
      "mono": "\"SF Mono\", \"Segoe UI Mono\", monospace"
    },
    "fontSize": {
      "h1": "2.5rem",
      "h2": "2rem",
      "h3": "1.5rem",
      "body": "1rem",
      "small": "0.875rem"
    }
  },
  "spacing": {
    "1": "0.25rem",
    "2": "0.5rem",
    "4": "1rem",
    "6": "1.5rem",
    "8": "2rem"
  }
}
```

## Design Review Checklist

Before implementing new features or components:

- [ ] Follows established color system
- [ ] Uses appropriate typography scale
- [ ] Maintains consistent spacing
- [ ] Implements responsive design
- [ ] Ensures accessibility standards
- [ ] Provides keyboard navigation
- [ ] Includes loading/error states
- [ ] Matches CLI terminology
- [ ] Optimized for performance
- [ ] Tested across browsers

## Future Enhancements

- [ ] Add theme customization support
- [ ] Implement dark/light mode toggle
- [ ] Create component library documentation
- [ ] Add animation library integration
- [ ] Design mobile-specific components
- [ ] Create design system documentation site
- [ ] Add Figma design tokens export
- [ ] Implement design token versioning

## References

- [GitHub Design System](https://primer.style/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Lucide Icons](https://lucide.dev/)
- [Web Content Accessibility Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
