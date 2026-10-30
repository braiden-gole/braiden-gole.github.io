class DiagramViewer {
    /**
     * Initializes a new instance of the DiagramViewer class.
     */
    constructor(viewportIdentifier, contentIdentifier) {
        // Retrieve and store reference to the primary viewport container element.
        this.viewportElement = document.getElementById(viewportIdentifier)
        
        // Retrieve and store reference to the inner SVG content wrapper element.
        this.contentElement = document.getElementById(contentIdentifier)
        
        // Initialize default view transform parameters for scale factor and pan offsets.
        this.currentScale = 1
        this.panX = 0
        this.panY = 0
        
        // Define incremental zoom step size and zoom level boundaries.
        this.scaleStep = 0.015
        this.minScale = 0.05
        this.maxScale = 20

        // Set initial flags and coordinates for tracking user drag and pinch gestures.
        this.isDragging = false
        this.isPinching = false
        this.startX = 0
        this.startY = 0
        this.initialPinchDistance = 0
        this.initialScaleBeforePinch = 1

        // Execute application startup logic and register interaction listeners.
        this.initializeViewerApplication()
    }

    /**
     * Initializes the viewer application state, parses URL query parameters, and loads the target SVG asset.
     */
    initializeViewerApplication() {
        // Parse incoming URL query string parameters to extract diagram options.
        const urlParameters = new URLSearchParams(window.location.search)
        
        // Extract target SVG file path and optional title override from URL parameters.
        const svgSourceFile = urlParameters.get('src') || urlParameters.get('file')
        const diagramTitle = urlParameters.get('title')

        // Locate title display element in the page header.
        const titleElement = document.getElementById('diagram-title')

        // Update document title and page header based on provided query parameters.
        if (titleElement) {
            if (diagramTitle) {
                // Decode URI components in the title parameter for clean text display.
                const decodedTitle = decodeURIComponent(diagramTitle)

                // Render decoded custom title in the header interface.
                titleElement.textContent = decodedTitle

                // Set browser tab title dynamically to match the diagram name.
                document.title = `${decodedTitle} | SVG Viewer`
            } else {
                // Fall back to displaying the raw filename or default title string.
                titleElement.textContent = svgSourceFile ? svgSourceFile.split('/').pop() : 'Untitled Diagram'
            }
        }

        // Trigger SVG document loading or display a missing asset warning.
        if (svgSourceFile) {
            // Asynchronously fetch and render the requested SVG diagram file.
            this.loadSvgDocument(svgSourceFile)
        } else if (titleElement) {
            // Display notice when no target SVG file is provided in URL parameters.
            titleElement.textContent = 'No SVG file specified'
        }

        // Register event handlers for mouse, touch, wheel, and button controls.
        this.attachInteractionEventListeners()
    }

    /**
     * Determines whether an SVG document is dark or light by examining its explicit background 
     * elements, node fill styles, or overall fill colors.
     */
    detectSvgIsDark(svgElement) {
        // 1. Check explicit background elements before removal.
        const bgElements = svgElement.querySelectorAll('rect.background, rect[width="100%"][height="100%"]')

        // Inspect background rectangles for explicit fill colors to evaluate luminance.
        for (const rect of bgElements) {
            // Read computed CSS fill property or standard fill attribute.
            const fill = window.getComputedStyle(rect).fill || rect.getAttribute('fill')

            // Evaluate darkness directly if a valid non-transparent fill color exists.
            if (fill && fill !== 'none' && fill !== 'transparent') {
                return this.isColorDark(fill)
            }
        }

        // 2. Sample fill colors of node shapes/paths across the SVG.
        const shapeElements = svgElement.querySelectorAll('rect, path, polygon, circle')

        // Track aggregate luminance values and total sampled elements.
        let totalLuminance = 0
        let sampledCount = 0

        // Iterate through a representative sample of shape elements across the SVG document.
        for (let index = 0; index < Math.min(shapeElements.length, 30); index++) {
            const element = shapeElements[index]
            const fill = window.getComputedStyle(element).fill || element.getAttribute('fill')

            // Calculate luminance for each non-transparent shape fill found.
            if (fill && fill !== 'none' && fill !== 'transparent') {
                const luminance = this.getColorLuminance(fill)

                // Accumulate valid luminance scores for average calculation.
                if (luminance !== null) {
                    totalLuminance += luminance
                    sampledCount++
                }
            }
        }

        // Return theme classification based on average sampled shape luminance.
        if (sampledCount > 0) {
            // If average sampled shape luminance is low (< 0.5), consider the diagram dark-themed.
            return (totalLuminance / sampledCount) < 0.5
        }

        // Default to light theme if undetectable.
        return false
    }

    /**
     * Helper to compute relative luminance (0 = black, 1 = white) from a CSS color string.
     */
    getColorLuminance(colorStr) {
        // Initialize red, green, and blue color component channels.
        let r = 0
        let g = 0
        let b = 0

        // Parse hex color formats including three-digit and six-digit syntax.
        if (colorStr.startsWith('#')) {
            let hex = colorStr.replace('#', '')
            
            // Expand shorthand three-character hex string to full six-character format.
            if (hex.length === 3) {
                hex = hex.split('').map(c => c + c).join('')
            }
            
            // Parse two-digit hex values into integer byte values.
            if (hex.length === 6) {
                r = parseInt(hex.substring(0, 2), 16)
                g = parseInt(hex.substring(2, 4), 16)
                b = parseInt(hex.substring(4, 6), 16)
            }
        } else if (colorStr.startsWith('rgb')) {
            // Extract numeric channel components from standard RGB string format.
            const matches = colorStr.match(/\d+/g)
            
            // Assign parsed RGB channel integers.
            if (matches && matches.length >= 3) {
                r = parseInt(matches[0], 10)
                g = parseInt(matches[1], 10)
                b = parseInt(matches[2], 10)
            }
        } else {
            // Return null if color string format cannot be parsed safely.
            return null
        }

        // Relative luminance formula (ITU-R BT.709).
        return (0.2126 * r + 0.7152 * g + 0.0722 * b) / 255
    }

    /**
     * Checks if a given color string represents a dark color.
     */
    isColorDark(colorStr) {
        // Compute numeric relative luminance score for the provided color.
        const luminance = this.getColorLuminance(colorStr)

        // Classify color as dark if calculated luminance falls below threshold.
        return luminance !== null ? luminance < 0.5 : false
    }

    /**
     * Switches the application theme (`theme-dark` / `theme-light` on `<html>`).
     */
    applyTheme(isDark) {
        // Obtain root HTML document element to adjust global CSS theme classes.
        const htmlElement = document.documentElement
        
        // Toggle dark and light theme selector classes on the root element.
        if (isDark) {
            // Remove light theme class and enforce dark mode styling globally.
            htmlElement.classList.remove('theme-light')
            htmlElement.classList.add('theme-dark')
        } else {
            // Remove dark theme class and enforce light mode styling globally.
            htmlElement.classList.remove('theme-dark')
            htmlElement.classList.add('theme-light')
        }
    }

    /**
     * Asynchronously fetches, cleans, scales, and injects an SVG document into the viewer container.
     */
    async loadSvgDocument(url) {
        // Abort loading sequence if target content container element is missing.
        if (!this.contentElement) {
            return
        }

        try {
            // Fetch remote SVG file asset over network.
            const response = await fetch(url)

            // Validate successful HTTP response status code before processing.
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }
            
            // Read raw text contents from HTTP response stream.
            let svgText = await response.text()

            // Strip whitespace around fetched SVG markup string.
            svgText = svgText.trim()

            // Inject raw SVG markup string into viewer container DOM node.
            this.contentElement.innerHTML = svgText

            // Query injected SVG root node for layout adjustments and styling.
            const svgElement = this.contentElement.querySelector('svg')

            if (svgElement) {
                // 1. Inject layout rules without overriding per-node text colors.
                const styleElement = document.createElementNS("http://www.w3.org/2000/svg", "style")                

                // Prepend injected style rules into the top of the SVG element.
                svgElement.prepend(styleElement)

                // Remove solid background rectangles obstructing full-canvas rendering.
                const bgElements = svgElement.querySelectorAll('rect.background, rect[width="100%"][height="100%"]')

                // Strip full-width background rectangles from the document structure.
                bgElements.forEach((rectangle) => rectangle.remove())

                // Set SVG presentation styles for clean rendering and overflow handling.
                svgElement.style.backgroundColor = 'transparent'
                svgElement.style.shapeRendering = 'geometricPrecision'
                svgElement.style.textRendering = 'geometricPrecision'
                svgElement.style.imageRendering = 'optimizeQuality'
                svgElement.style.overflow = 'visible'

                // Reset rigid dimensional attributes in favor of viewBox scaling.
                svgElement.removeAttribute('width')
                svgElement.removeAttribute('height')
                svgElement.setAttribute('preserveAspectRatio', 'xMidYMid meet')

                // 2. Safely expand viewBox using the native SVG bounding box.
                try {
                    // Calculate native geometric bounding box of the parsed SVG content.
                    const bbox = svgElement.getBBox()
                    
                    // Verify that calculated bounding box dimensions are valid positive numbers.
                    if (bbox && bbox.width > 0 && bbox.height > 0) {
                        // Compute dynamic canvas margin padding based on overall diagram dimensions.
                        const padding = Math.max(60, Math.max(bbox.width, bbox.height) * 0.08)
                        const naturalWidth = Math.ceil(bbox.width + padding * 2)
                        const naturalHeight = Math.ceil(bbox.height + padding * 2)

                        // Update SVG viewBox attribute with padded dimensions.
                        svgElement.setAttribute(
                            'viewBox',
                            `${bbox.x - padding} ${bbox.y - padding} ${naturalWidth} ${naturalHeight}`
                        )

                        // Set content dimensions to natural width and height.
                        this.contentElement.style.width = `${naturalWidth}px`
                        this.contentElement.style.height = `${naturalHeight}px`
                        
                        svgElement.style.width = '100%'
                        svgElement.style.height = '100%'
                    }
                } catch (error) {
                    // Log warning if native bounding box calculation fails.
                    console.warn("Could not calculate SVG bounding box:", error)
                }

                // Schedule initial viewport zoom and pan centering alignment.
                requestAnimationFrame(() => {
                    this.resetZoomAndPan()
                })
            }
        } catch (error) {
            // Display visual error message inside container if loading fails.
            if (this.contentElement) {
                this.contentElement.innerHTML = `<p class="has-text-danger p-4">Error loading SVG file: ${error.message}</p>`
            }
        }
    }

    /**
     * Calculates the Euclidean distance between two active touch points for pinch-to-zoom operations.
     */
    getPinchDistance(touches) {
        // Calculate horizontal and vertical point deltas between touch inputs.
        const deltaX = touches[0].clientX - touches[1].clientX
        const deltaY = touches[0].clientY - touches[1].clientY

        // Return straight-line Euclidean distance using the Pythagorean theorem.
        return Math.sqrt(deltaX * deltaX + deltaY * deltaY)
    }

    /**
     * Attaches user interface, mouse, touch, scroll, and toolbar event listeners to the viewer environment.
     */
    attachInteractionEventListeners() {
        // Locate navigation back button element in document DOM.
        const backButton = document.getElementById('back-btn') || document.querySelector('.back-btn')

        // Register navigation click listener for return navigation button.
        if (backButton) {
            backButton.addEventListener('click', (event) => {
                // Prevent default anchor navigation behavior.
                event.preventDefault()
                
                // Extract target navigation selection from active query parameters.
                const urlParameters = new URLSearchParams(window.location.search)
                const targetSelection = urlParameters.get('target') || urlParameters.get("from") || "diagrams"

                // Use query parameters instead of a hash to prevent premature page jumping.
                window.location.href = `index.html?target=${targetSelection}`
            })
        }

        // Abort event listener registration if viewport container is unavailable.
        if (!this.viewportElement) {
            return
        }

        // Handle mouse press event to initiate dragging interaction.
        this.viewportElement.addEventListener('mousedown', (event) => {
            // Restrict dragging actions exclusively to primary left mouse button presses.
            if (event.button !== 0) {
                return
            }

            // Enable dragging active flag state.
            this.isDragging = true

            // Record initial drag offset coordinates relative to active pan position.
            this.startX = event.clientX - this.panX
            this.startY = event.clientY - this.panY

            // Update mouse cursor indicator to grabbing state.
            this.viewportElement.style.cursor = 'grabbing'
        })

        // Handle global mouse movement events during dragging interactions.
        window.addEventListener('mousemove', (event) => {
            // Ignore mouse movement if user is not actively dragging.
            if (!this.isDragging) {
                return
            }

            // Prevent default drag selection behavior across window.
            event.preventDefault()

            // Update active pan offset coordinates based on cursor motion.
            this.panX = event.clientX - this.startX
            this.panY = event.clientY - this.startY

            // Apply updated transform matrix to content container.
            this.updateTransformMatrix()
        })

        // Handle mouse release events to conclude drag operations.
        window.addEventListener('mouseup', () => {
            // Reset dragging state and cursor styles upon mouse button release.
            if (this.isDragging) {
                this.isDragging = false

                // Restore grab cursor style on viewport container.
                if (this.viewportElement) {
                    this.viewportElement.style.cursor = 'grab'
                }
            }
        })

        // Reset drag states if cursor moves outside viewport boundary.
        this.viewportElement.addEventListener('mouseleave', () => {
            // Cancel active drag gesture on cursor exit.
            if (this.isDragging) {
                this.isDragging = false
                this.viewportElement.style.cursor = 'grab'
            }
        })

        // Register touch start listener for single-finger drag and multi-finger pinch gestures.
        this.viewportElement.addEventListener('touchstart', (event) => {
            // Differentiate between single-touch dragging and dual-touch pinch zooming.
            if (event.touches.length === 1) {
                // Initialize single-finger drag gesture tracking parameters.
                this.isPinching = false
                this.isDragging = true
                this.startX = event.touches[0].clientX - this.panX
                this.startY = event.touches[0].clientY - this.panY
            } else if (event.touches.length === 2) {
                // Initialize dual-finger pinch gesture tracking parameters.
                this.isDragging = false
                this.isPinching = true
                this.initialPinchDistance = this.getPinchDistance(event.touches)
                this.initialScaleBeforePinch = this.currentScale
            }
        }, { passive: true })

        // Process touch movement updates for active touch gestures.
        window.addEventListener('touchmove', (event) => {
            // Process dual-touch pinch zoom scale updates.
            if (event.touches.length === 2 && this.isPinching) {
                // Prevent standard browser page scaling during touch move.
                event.preventDefault()

                // Calculate current distance between dual touch contacts.
                const currentDistance = this.getPinchDistance(event.touches)

                // Compute dynamic scale factor relative to initial touch distance.
                if (this.initialPinchDistance > 0) {
                    const scaleFactor = currentDistance / this.initialPinchDistance

                    // Clamp new scale factor within defined minimum and maximum limits.
                    this.currentScale = Math.max(
                        this.minScale,
                        Math.min(this.maxScale, this.initialScaleBeforePinch * scaleFactor)
                    )

                    // Render updated transformation scale.
                    this.updateTransformMatrix()
                }
            } else if (this.isDragging && event.touches.length === 1) {
                // Update pan coordinates according to single-finger touch position.
                this.panX = event.touches[0].clientX - this.startX
                this.panY = event.touches[0].clientY - this.startY
                
                // Render updated transformation translation.
                this.updateTransformMatrix()
            }
        }, { passive: false })

        // Process touch release events to conclude touch gestures cleanly.
        window.addEventListener('touchend', (event) => {
            // Discontinue pinch zoom mode when fewer than two touch contacts remain.
            if (event.touches.length < 2) {
                this.isPinching = false
            }

            // Discontinue or transition drag mode depending on remaining active touches.
            if (event.touches.length === 0) {
                // End dragging when all touch points leave screen.
                this.isDragging = false
            } else if (event.touches.length === 1) {
                // Transition back to single-finger dragging when one touch remains.
                this.isDragging = true
                this.startX = event.touches[0].clientX - this.panX
                this.startY = event.touches[0].clientY - this.panY
            }
        })

        // Register mouse wheel event listener for zooming controls.
        this.viewportElement.addEventListener('wheel', (event) => {
            // Prevent default page scrolling behavior on wheel action.
            event.preventDefault()

            // Determine zoom direction based on mouse wheel delta orientation.
            const zoomDirection = event.deltaY < 0 ? 1 : -1
            const zoomFactor = 1 + (zoomDirection * this.scaleStep)

            // Calculate and clamp new zoom scale within allowed boundary range.
            this.currentScale = Math.max(this.minScale, Math.min(this.maxScale, this.currentScale * zoomFactor))
            
            // Render updated scale transformation matrix.
            this.updateTransformMatrix()
        }, { passive: false })

        // Query control toolbar button elements from the document DOM.
        const zoomInButton = document.getElementById('zoom-in-btn')
        const zoomOutButton = document.getElementById('zoom-out-btn')
        const resetZoomButton = document.getElementById('reset-zoom-btn')

        // Attach click handler for zoom-in control button.
        if (zoomInButton) {
            zoomInButton.addEventListener('click', () => this.zoomIn())
        }

        // Attach click handler for zoom-out control button.
        if (zoomOutButton) {
            zoomOutButton.addEventListener('click', () => this.zoomOut())
        }

        // Attach click handler for zoom reset control button.
        if (resetZoomButton) {
            resetZoomButton.addEventListener('click', () => this.resetZoomAndPan())
        }

        // Register window resize listener to auto-fit view on window size changes.
        window.addEventListener('resize', () => this.resetZoomAndPan())
        
        // Set initial default grabbing cursor on viewport container element.
        this.viewportElement.style.cursor = 'grab'
    }

    /**
     * Increases the current viewing scale factor incrementally.
     */
    zoomIn() {
        // Increment scale factor while enforcing upper maximum zoom boundary.
        this.currentScale = Math.min(this.maxScale, this.currentScale + 0.1)
        
        // Apply updated transformation matrix to content container.
        this.updateTransformMatrix()
    }

    /**
     * Decreases the current viewing scale factor incrementally.
     */
    zoomOut() {
        // Decrement scale factor while enforcing lower minimum zoom boundary.
        this.currentScale = Math.max(this.minScale, this.currentScale - 0.1)
        
        // Apply updated transformation matrix to content container.
        this.updateTransformMatrix()
    }

    /**
     * Resets pan coordinates and computes the optimal scale to fit the SVG centered within the viewport.
     */
    resetZoomAndPan() {
        // Reset translation offset values back to origin coordinates.
        this.panX = 0
        this.panY = 0

        if (!this.contentElement) {
            return
        }

        const svgElement = this.contentElement.querySelector('svg')

        if (svgElement && this.viewportElement) {
            const viewportWidth = this.viewportElement.clientWidth || window.innerWidth
            const viewportHeight = this.viewportElement.clientHeight || window.innerHeight

            let svgWidth = parseFloat(this.contentElement.style.width) || 800
            let svgHeight = parseFloat(this.contentElement.style.height) || 600

            // Padding ratio to ensure proper spacing on viewport boundaries.
            const paddingRatio = 0.95
            const scaleX = (viewportWidth * paddingRatio) / svgWidth
            const scaleY = (viewportHeight * paddingRatio) / svgHeight

            // Pick scale factor to fit within both dimensions.
            let computedScale = Math.min(scaleX, scaleY)

            this.currentScale = Math.max(this.minScale, Math.min(this.maxScale, computedScale))
        } else {
            this.currentScale = 1
        }

        this.updateTransformMatrix()
    }

    /**
     * Applies the current pan translations and scale factor to the CSS transform matrix of the content container.
     */
    updateTransformMatrix() {
        // Render current translation offsets and scale factor as CSS transform rule centered relative to the viewport.
        if (this.contentElement) {
            this.contentElement.style.transform = `translate(calc(-50% + ${this.panX}px), calc(-50% + ${this.panY}px)) scale(${this.currentScale})`
            this.contentElement.style.transformOrigin = 'center center'
        }
    }
}

/**
 * Instantiates the DiagramViewer class once the DOM content has fully loaded.
 */
document.addEventListener("DOMContentLoaded", () => {
    // Create global instance of DiagramViewer bound to viewer HTML containers.
    window.diagramViewerInstance = new DiagramViewer("svg-viewport", "svg-content")
})