class DiagramViewer {
    /**
     * Initializes a new instance of the DiagramViewer class.
     */
    constructor(viewportIdentifier, contentIdentifier) {
        this.viewportElement = document.getElementById(viewportIdentifier)
        this.contentElement = document.getElementById(contentIdentifier)
        
        this.currentScale = 1
        this.panX = 0
        this.panY = 0
        
        this.scaleStep = 0.015
        this.minScale = 0.05
        this.maxScale = 20

        this.isDragging = false
        this.isPinching = false
        this.startX = 0
        this.startY = 0
        this.initialPinchDistance = 0
        this.initialScaleBeforePinch = 1

        this.initializeViewerApplication()
    }

    /**
     * Initializes the viewer application state, parses URL query parameters, and loads the target SVG asset.
     */
    initializeViewerApplication() {
        const urlParameters = new URLSearchParams(window.location.search)
        const svgSourceFile = urlParameters.get('src') || urlParameters.get('file')
        const diagramTitle = urlParameters.get('title')

        const titleElement = document.getElementById('diagram-title')

        if (titleElement) {
            if (diagramTitle) {
                const decodedTitle = decodeURIComponent(diagramTitle)

                titleElement.textContent = decodedTitle

                document.title = `${decodedTitle} | SVG Viewer`
            } else {
                titleElement.textContent = svgSourceFile ? svgSourceFile.split('/').pop() : 'Untitled Diagram'
            }
        }

        if (svgSourceFile) {
            this.loadSvgDocument(svgSourceFile)
        } else if (titleElement) {
            titleElement.textContent = 'No SVG file specified'
        }

        this.attachInteractionEventListeners()
    }

    /**
     * Asynchronously fetches, cleans, scales, and injects an SVG document into the viewer container.
     */
    async loadSvgDocument(url) {
        if (!this.contentElement) {
            return
        }

        try {
            const response = await fetch(url)

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`)
            }
            
            let svgText = await response.text()
            
            svgText = svgText.trim()

            this.contentElement.innerHTML = svgText

            const svgElement = this.contentElement.querySelector('svg')

            if (svgElement) {
                // 1. Inject global SVG rules to guarantee text visibility and overflow padding.
                const styleElement = document.createElementNS("http://www.w3.org/2000/svg", "style")

                styleElement.textContent = `
                    svg { overflow: visible !important; }
                    g.node foreignObject { 
                        overflow: visible !important; 
                    }
                    g.node foreignObject > div { 
                        white-space: normal !important; 
                        word-break: break-word !important; 
                        box-sizing: border-box !important;
                        overflow: visible !important;
                        display: flex !important;
                        align-items: center !important;
                        justify-content: center !important;
                        height: 100% !important;
                        width: 100% !important;
                        color: #ffffff !important;
                    }
                    g.node span.nodeLabel, g.node span.nodeLabel p, g.node p, g.node label { 
                        color: #ffffff !important;
                        fill: #ffffff !important;
                        font-family: system-ui, -apple-system, sans-serif !important; 
                        font-size: 13px !important;
                        line-height: 1.4 !important;
                        margin: 0 !important;
                        text-align: center !important;
                    }
                `
                svgElement.prepend(styleElement)

                // Remove solid background rectangles obstructing full-canvas rendering.
                const bgElements = svgElement.querySelectorAll('rect.background, rect[width="100%"][height="100%"]')

                bgElements.forEach((rectangle) => rectangle.remove())

                svgElement.style.backgroundColor = 'transparent'
                svgElement.style.shapeRendering = 'geometricPrecision'
                svgElement.style.textRendering = 'geometricPrecision'
                svgElement.style.imageRendering = 'optimizeQuality'
                svgElement.style.overflow = 'visible'

                svgElement.removeAttribute('width')
                svgElement.removeAttribute('height')
                svgElement.setAttribute('preserveAspectRatio', 'xMidYMid meet')

                // 2. Safely expand viewBox using the native SVG bounding box.
                try {
                    const bbox = svgElement.getBBox()
                    
                    if (bbox && bbox.width > 0 && bbox.height > 0) {
                        const padding = Math.max(60, Math.max(bbox.width, bbox.height) * 0.08)
                        const naturalWidth = Math.ceil(bbox.width + padding * 2)
                        const naturalHeight = Math.ceil(bbox.height + padding * 2)

                        svgElement.setAttribute(
                            'viewBox',
                            `${bbox.x - padding} ${bbox.y - padding} ${naturalWidth} ${naturalHeight}`
                        )

                        const HD_DENSITY_FACTOR = 2
                        const renderWidth = naturalWidth * HD_DENSITY_FACTOR
                        const renderHeight = naturalHeight * HD_DENSITY_FACTOR

                        svgElement.style.width = `${renderWidth}px`
                        svgElement.style.height = `${renderHeight}px`
                        this.contentElement.style.width = `${renderWidth}px`
                        this.contentElement.style.height = `${renderHeight}px`
                    }
                } catch (error) {
                    console.warn("Could not calculate SVG bounding box:", error)
                }

                requestAnimationFrame(() => {
                    this.resetZoomAndPan()
                })
            }
        } catch (error) {
            if (this.contentElement) {
                this.contentElement.innerHTML = `<p class="has-text-danger p-4">Error loading SVG file: ${error.message}</p>`
            }
        }
    }

    /**
     * Calculates the Euclidean distance between two active touch points for pinch-to-zoom operations.
     */
    getPinchDistance(touches) {
        const deltaX = touches[0].clientX - touches[1].clientX
        const deltaY = touches[0].clientY - touches[1].clientY

        return Math.sqrt(deltaX * deltaX + deltaY * deltaY)
    }

    /**
     * Attaches user interface, mouse, touch, scroll, and toolbar event listeners to the viewer environment.
     */
    attachInteractionEventListeners() {
        const backButton = document.getElementById('back-btn') || document.querySelector('.back-btn')

        if (backButton) {
            backButton.addEventListener('click', (event) => {
                event.preventDefault()
                
                const urlParameters = new URLSearchParams(window.location.search)
                const targetSelection = urlParameters.get('target') || urlParameters.get("from") || "diagrams"

                // Use query parameters instead of a hash to prevent premature page jumping.
                window.location.href = `index.html?target=${targetSelection}`
            })
        }

        if (!this.viewportElement) {
            return
        }

        this.viewportElement.addEventListener('mousedown', (event) => {
            if (event.button !== 0) {
                return
            }

            this.isDragging = true

            this.startX = event.clientX - this.panX
            this.startY = event.clientY - this.panY

            this.viewportElement.style.cursor = 'grabbing'
        })

        window.addEventListener('mousemove', (event) => {
            if (!this.isDragging) {
                return
            }

            event.preventDefault()

            this.panX = event.clientX - this.startX
            this.panY = event.clientY - this.startY

            this.updateTransformMatrix()
        })

        window.addEventListener('mouseup', () => {
            if (this.isDragging) {
                this.isDragging = false

                if (this.viewportElement) {
                    this.viewportElement.style.cursor = 'grab'
                }
            }
        })

        this.viewportElement.addEventListener('mouseleave', () => {
            if (this.isDragging) {
                this.isDragging = false
                this.viewportElement.style.cursor = 'grab'
            }
        })

        this.viewportElement.addEventListener('touchstart', (event) => {
            if (event.touches.length === 1) {
                this.isPinching = false
                this.isDragging = true
                this.startX = event.touches[0].clientX - this.panX
                this.startY = event.touches[0].clientY - this.panY
            } else if (event.touches.length === 2) {
                this.isDragging = false
                this.isPinching = true
                this.initialPinchDistance = this.getPinchDistance(event.touches)
                this.initialScaleBeforePinch = this.currentScale
            }
        }, { passive: true })

        window.addEventListener('touchmove', (event) => {
            if (event.touches.length === 2 && this.isPinching) {
                event.preventDefault()

                const currentDistance = this.getPinchDistance(event.touches)

                if (this.initialPinchDistance > 0) {
                    const scaleFactor = currentDistance / this.initialPinchDistance

                    this.currentScale = Math.max(
                        this.minScale,
                        Math.min(this.maxScale, this.initialScaleBeforePinch * scaleFactor)
                    )

                    this.updateTransformMatrix()
                }
            } else if (this.isDragging && event.touches.length === 1) {
                this.panX = event.touches[0].clientX - this.startX
                this.panY = event.touches[0].clientY - this.startY
                this.updateTransformMatrix()
            }
        }, { passive: false })

        window.addEventListener('touchend', (event) => {
            if (event.touches.length < 2) {
                this.isPinching = false
            }

            if (event.touches.length === 0) {
                this.isDragging = false
            } else if (event.touches.length === 1) {
                this.isDragging = true
                this.startX = event.touches[0].clientX - this.panX
                this.startY = event.touches[0].clientY - this.panY
            }
        })

        this.viewportElement.addEventListener('wheel', (event) => {
            event.preventDefault()

            const zoomDirection = event.deltaY < 0 ? 1 : -1
            const zoomFactor = 1 + (zoomDirection * this.scaleStep)

            this.currentScale = Math.max(this.minScale, Math.min(this.maxScale, this.currentScale * zoomFactor))
            this.updateTransformMatrix()
        }, { passive: false })

        const zoomInButton = document.getElementById('zoom-in-btn')
        const zoomOutButton = document.getElementById('zoom-out-btn')
        const resetZoomButton = document.getElementById('reset-zoom-btn')

        if (zoomInButton) {
            zoomInButton.addEventListener('click', () => this.zoomIn())
        }

        if (zoomOutButton) {
            zoomOutButton.addEventListener('click', () => this.zoomOut())
        }

        if (resetZoomButton) {
            resetZoomButton.addEventListener('click', () => this.resetZoomAndPan())
        }

        window.addEventListener('resize', () => this.resetZoomAndPan())
        
        this.viewportElement.style.cursor = 'grab'
    }

    /**
     * Increases the current viewing scale factor incrementally.
     */
    zoomIn() {
        this.currentScale = Math.min(this.maxScale, this.currentScale + 0.1)
        this.updateTransformMatrix()
    }

    /**
     * Decreases the current viewing scale factor incrementally.
     */
    zoomOut() {
        this.currentScale = Math.max(this.minScale, this.currentScale - 0.1)
        this.updateTransformMatrix()
    }

    /**
     * Resets pan coordinates and computes the optimal scale to fit the SVG within the viewport.
     */
    resetZoomAndPan() {
        this.panX = 0
        this.panY = 0

        if (!this.contentElement) {
            return
        }

        const svgElement = this.contentElement.querySelector('svg')

        if (svgElement && this.viewportElement) {
            const viewportWidth = this.viewportElement.clientWidth || window.innerWidth
            const viewportHeight = this.viewportElement.clientHeight || window.innerHeight

            let svgWidth = parseFloat(svgElement.style.width) || 800
            let svgHeight = parseFloat(svgElement.style.height) || 600

            const paddingRatio = 0.85
            const scaleX = (viewportWidth * paddingRatio) / svgWidth
            const scaleY = (viewportHeight * paddingRatio) / svgHeight

            this.currentScale = Math.min(scaleX, scaleY)
            this.currentScale = Math.max(this.minScale, Math.min(this.maxScale, this.currentScale))
        } else {
            this.currentScale = 1
        }

        this.updateTransformMatrix()
    }

    /**
     * Applies the current pan translations and scale factor to the CSS transform matrix of the content container.
     */
    updateTransformMatrix() {
        if (this.contentElement) {
            this.contentElement.style.transform = `translate(calc(-50% + ${this.panX}px), calc(-50% + ${this.panY}px)) scale(${this.currentScale})`
        }
    }
}

/**
 * Instantiates the DiagramViewer class once the DOM content has fully loaded.
 */
document.addEventListener("DOMContentLoaded", () => {
    window.diagramViewerInstance = new DiagramViewer("svg-viewport", "svg-content")
})