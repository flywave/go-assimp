# Open Asset Import Library (assimp)
# ----------------------------------------------------------------------
# Copyright (c) 2006-2022, assimp team
#
# All rights reserved.
#
# Redistribution and use of this software in source and binary forms,
# with or without modification, are permitted provided that the
# following conditions are met:
#
# * Redistributions of source code must retain the above
#   copyright notice, this list of conditions and the
#   following disclaimer.
#
# * Redistributions in binary form must reproduce the above
#   copyright notice, this list of conditions and the
#   following disclaimer in the documentation and/or other
#   materials provided with the distribution.
#
# * Neither the name of the assimp team, nor the names of its
#   contributors may be used to endorse or promote products
#   derived from this software without specific prior
#   written permission of the assimp team.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
# "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
# LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
# A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
# OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
# SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
# LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
# DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
# THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
# (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
#----------------------------------------------------------------------
SET(CMAKE_POLICY_DEFAULT_CMP0012 NEW)
SET(CMAKE_POLICY_DEFAULT_CMP0074 NEW)
SET(CMAKE_POLICY_DEFAULT_CMP0092 NEW)

CMAKE_MINIMUM_REQUIRED( VERSION 3.10 )

# Disabled importers: m3d for 5.1
ADD_DEFINITIONS( -DASSIMP_BUILD_NO_M3D_IMPORTER)
ADD_DEFINITIONS( -DASSIMP_BUILD_NO_M3D_EXPORTER)
# Toggles the use of the hunter package manager
option(ASSIMP_HUNTER_ENABLED "Enable Hunter package manager support" OFF)

IF(ASSIMP_HUNTER_ENABLED)
  include("cmake-modules/HunterGate.cmake")
  HunterGate(
    URL "https://github.com/cpp-pm/hunter/archive/v0.23.311.tar.gz"
    SHA1 "1a82b9b73055879181cb1466b2ab5d48ee8ae410"
  )

  add_definitions(-DASSIMP_USE_HUNTER)
ENDIF()

PROJECT(Assimp VERSION 5.1.6)

# All supported options ###############################################

OPTION( BUILD_SHARED_LIBS
  "Build package with shared libraries."
  OFF
)
OPTION( ASSIMP_BUILD_FRAMEWORK
  "Build package as Mac OS X Framework bundle."
  OFF
)
OPTION( ASSIMP_DOUBLE_PRECISION
  "Set to ON to enable double precision processing"
  OFF
)
OPTION( ASSIMP_OPT_BUILD_PACKAGES
  "Set to ON to generate CPack configuration files and packaging targets"
  OFF
)
OPTION( ASSIMP_ANDROID_JNIIOSYSTEM
  "Android JNI IOSystem support is active"
  OFF
)
OPTION( ASSIMP_NO_EXPORT
  "Disable Assimp's export functionality."
  OFF
)
OPTION( ASSIMP_BUILD_ZLIB
  "Build your own zlib"
  OFF
)
OPTION( ASSIMP_BUILD_ASSIMP_TOOLS
  "If the supplementary tools for Assimp are built in addition to the library."
  OFF
)
OPTION ( ASSIMP_BUILD_SAMPLES
  "If the official samples are built as well (needs Glut)."
  OFF
)
OPTION ( ASSIMP_BUILD_TESTS
  "If the test suite for Assimp is built in addition to the library."
  OFF
)
OPTION ( ASSIMP_COVERALLS
  "Enable this to measure test coverage."
  OFF
)
OPTION( ASSIMP_INSTALL
  "Disable this if you want to use assimp as a submodule."
  OFF
)
OPTION ( ASSIMP_ERROR_MAX
  "Enable all warnings."
  OFF
)
OPTION ( ASSIMP_ASAN
  "Enable AddressSanitizer."
  OFF
)
OPTION ( ASSIMP_UBSAN
  "Enable Undefined Behavior sanitizer."
  OFF
)
OPTION ( ASSIMP_BUILD_DOCS
  "Build documentation using Doxygen."
  OFF
)
OPTION( ASSIMP_INJECT_DEBUG_POSTFIX
  "Inject debug postfix in .a/.so/.dll lib names"
  OFF
)

OPTION ( ASSIMP_IGNORE_GIT_HASH
   "Don't call git to get the hash."
   OFF
)

IF ( WIN32 )
  # Use subset of Windows.h
  ADD_DEFINITIONS( -DWIN32_LEAN_AND_MEAN )

  IF(MSVC)
    OPTION ( ASSIMP_BUILD_ASSIMP_VIEW
      "If the Assimp view tool is built. (requires DirectX)"
      OFF )

    OPTION( ASSIMP_INSTALL_PDB
      "Install MSVC debug files."
      ON )
    IF(NOT (MSVC_VERSION LESS 1900))
      # Multibyte character set is deprecated since at least MSVC2015 (possibly earlier)
      ADD_DEFINITIONS( -DUNICODE -D_UNICODE )
    ENDIF()
  ENDIF()
ENDIF()

IF (IOS AND NOT ASSIMP_HUNTER_ENABLED)
  IF (NOT CMAKE_BUILD_TYPE)
    SET(CMAKE_BUILD_TYPE "Release")
  ENDIF ()
  ADD_DEFINITIONS(-DENABLE_BITCODE)
ENDIF ()


IF (ASSIMP_BUILD_FRAMEWORK)
  SET (BUILD_SHARED_LIBS ON)
  MESSAGE(STATUS "Framework bundle building enabled")
ENDIF()

IF(NOT BUILD_SHARED_LIBS)
  MESSAGE(STATUS "Shared libraries disabled")
  SET(LINK_SEARCH_START_STATIC TRUE)
ELSE()
  MESSAGE(STATUS "Shared libraries enabled")
ENDIF()

# Define here the needed parameters
SET (ASSIMP_VERSION_MAJOR ${PROJECT_VERSION_MAJOR})
SET (ASSIMP_VERSION_MINOR ${PROJECT_VERSION_MINOR})
SET (ASSIMP_VERSION_PATCH ${PROJECT_VERSION_PATCH})
SET (ASSIMP_VERSION ${ASSIMP_VERSION_MAJOR}.${ASSIMP_VERSION_MINOR}.${ASSIMP_VERSION_PATCH})
SET (ASSIMP_SOVERSION 5)

SET( ASSIMP_PACKAGE_VERSION "0" CACHE STRING "the package-specific version used for uploading the sources" )
if(NOT ASSIMP_HUNTER_ENABLED)
  # Enable C++11 support globally
  set(CMAKE_CXX_STANDARD 11)
  set(CMAKE_CXX_STANDARD_REQUIRED ON)
  set(CMAKE_C_STANDARD 99)
endif()

IF(NOT ASSIMP_IGNORE_GIT_HASH)
  # Get the current working branch
  EXECUTE_PROCESS(
    COMMAND git rev-parse --abbrev-ref HEAD
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    OUTPUT_VARIABLE GIT_BRANCH
    OUTPUT_STRIP_TRAILING_WHITESPACE
    ERROR_QUIET
  )

  # Get the latest abbreviated commit hash of the working branch
  EXECUTE_PROCESS(
    COMMAND git rev-parse --short=8 HEAD
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    OUTPUT_VARIABLE GIT_COMMIT_HASH
    OUTPUT_STRIP_TRAILING_WHITESPACE
    ERROR_QUIET
  )
ENDIF()

IF(NOT GIT_COMMIT_HASH)
  SET(GIT_COMMIT_HASH 0)
ENDIF()

IF(ASSIMP_DOUBLE_PRECISION)
  ADD_DEFINITIONS(-DASSIMP_DOUBLE_PRECISION)
ENDIF()

CONFIGURE_FILE(
  ${CMAKE_CURRENT_LIST_DIR}/revision.h.in
  ${CMAKE_CURRENT_BINARY_DIR}/revision.h
)

CONFIGURE_FILE(
  ${CMAKE_CURRENT_LIST_DIR}/include/assimp/config.h.in
  ${CMAKE_CURRENT_BINARY_DIR}/include/assimp/config.h
)

INCLUDE_DIRECTORIES( BEFORE
  ./
  code/
  include
  ${CMAKE_CURRENT_BINARY_DIR}
  ${CMAKE_CURRENT_BINARY_DIR}/include
)

LIST(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake-modules" )
SET(LIBASSIMP_COMPONENT     "libassimp${ASSIMP_VERSION_MAJOR}.${ASSIMP_VERSION_MINOR}.${ASSIMP_VERSION_PATCH}" )
SET(LIBASSIMP-DEV_COMPONENT "libassimp${ASSIMP_VERSION_MAJOR}.${ASSIMP_VERSION_MINOR}.${ASSIMP_VERSION_PATCH}-dev" )
SET(CPACK_COMPONENTS_ALL assimp-bin ${LIBASSIMP_COMPONENT} ${LIBASSIMP-DEV_COMPONENT} assimp-dev)
SET(ASSIMP_LIBRARY_SUFFIX "" CACHE STRING "Suffix to append to library names")

IF( UNIX )
  # Use GNUInstallDirs for Unix predefined directories
  INCLUDE(GNUInstallDirs)
ENDIF()

# Grouped compiler settings ########################################
IF ((CMAKE_C_COMPILER_ID MATCHES "GNU") AND NOT MINGW)
  IF(NOT ASSIMP_HUNTER_ENABLED)
    SET(CMAKE_CXX_STANDARD 11)
    SET(CMAKE_POSITION_INDEPENDENT_CODE ON)
  ENDIF()
  # hide all not-exported symbols
  IF(CMAKE_HOST_SYSTEM_PROCESSOR MATCHES "mips64" )
	SET(CMAKE_CXX_FLAGS "-mxgot -fvisibility=hidden -fno-strict-aliasing -Wall ${CMAKE_CXX_FLAGS}")
	SET(CMAKE_C_FLAGS "-fno-strict-aliasing ${CMAKE_C_FLAGS}")
	SET(LIBSTDC++_LIBRARIES -lstdc++)
  ELSE()
	SET(CMAKE_CXX_FLAGS "-fvisibility=hidden -fno-strict-aliasing -Wall ${CMAKE_CXX_FLAGS}")
	SET(CMAKE_C_FLAGS "-fno-strict-aliasing ${CMAKE_C_FLAGS}")
	SET(LIBSTDC++_LIBRARIES -lstdc++)
  ENDIF()
ELSEIF(MSVC)
  # enable multi-core compilation with MSVC
  IF(CMAKE_CXX_COMPILER_ID MATCHES "Clang" ) # clang-cl
    ADD_COMPILE_OPTIONS(/bigobj /W4 /WX )
  ELSE() # msvc
    ADD_COMPILE_OPTIONS(/MP /bigobj /W4 /WX)
  ENDIF()
  # disable "elements of array '' will be default initialized" warning on MSVC2013
  IF(MSVC12)
    ADD_COMPILE_OPTIONS(/wd4351)	
  ENDIF()
  ADD_COMPILE_OPTIONS(/wd4244) #supress warning for double to float conversion if Double precission is activated
  SET(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} /D_DEBUG /Zi /Od")
  SET(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} /Zi")
  SET(CMAKE_SHARED_LINKER_FLAGS_RELEASE "${CMAKE_SHARED_LINKER_FLAGS_RELEASE} /DEBUG:FULL /PDBALTPATH:%_PDB% /OPT:REF /OPT:ICF")
ELSEIF (CMAKE_CXX_COMPILER_ID MATCHES "Clang" )
  IF(NOT ASSIMP_HUNTER_ENABLED)
    SET(CMAKE_CXX_STANDARD 11)
    SET(CMAKE_POSITION_INDEPENDENT_CODE ON)
  ENDIF()
  SET(CMAKE_CXX_FLAGS "-fvisibility=hidden -fno-strict-aliasing -Wall -Wno-long-long  -Wno-shorten-64-to-32 -Wno-deprecated-declarations -Wno-unused-but-set-variable ${CMAKE_CXX_FLAGS}" )
  SET(CMAKE_C_FLAGS "-fno-strict-aliasing  -Wno-shorten-64-to-32 -Wno-unused-but-set-variable ${CMAKE_C_FLAGS}")
ELSEIF( MINGW )
  IF (CMAKE_CXX_COMPILER_VERSION VERSION_LESS 7.0)
    message(FATAL_ERROR "MinGW is too old to be supported. Please update MinGW and try again.")
  ELSEIF(CMAKE_CXX_COMPILER_VERSION VERSION_LESS 7.3)
    message(WARNING "MinGW is old, if you experience errors, update MinGW.")
  ENDIF()
  IF(NOT ASSIMP_HUNTER_ENABLED)
    SET(CMAKE_CXX_FLAGS "-std=gnu++11 ${CMAKE_CXX_FLAGS}")
    SET(CMAKE_C_FLAGS "-fPIC ${CMAKE_C_FLAGS}")
  ENDIF()
  SET(CMAKE_CXX_FLAGS "-fvisibility=hidden -fno-strict-aliasing -Wall -Wno-long-long -Wa,-mbig-obj -O3 ${CMAKE_CXX_FLAGS}")
  SET(CMAKE_C_FLAGS "-fno-strict-aliasing ${CMAKE_C_FLAGS}")
ENDIF()

IF ( IOS AND NOT ASSIMP_HUNTER_ENABLED)
  IF (CMAKE_BUILD_TYPE STREQUAL "Debug")
    SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fembed-bitcode -Og")
    SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fembed-bitcode -Og")
  ELSE()
    SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fembed-bitcode -O3")
    SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fembed-bitcode -O3")
    # Experimental for pdb generation
  ENDIF()
ENDIF()

IF (ASSIMP_COVERALLS)
  MESSAGE(STATUS "Coveralls enabled")
  INCLUDE(Coveralls)
  SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -g -O0 -fprofile-arcs -ftest-coverage")
  SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g -O0 -fprofile-arcs -ftest-coverage")
ENDIF()

IF (ASSIMP_ERROR_MAX)
  MESSAGE(STATUS "Turning on all warnings")
  IF (MSVC)
    ADD_COMPILE_OPTIONS(/W4) # NB: there is a /Wall option, pedantic mode
  ELSE()
    SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall")
    SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wall")
  ENDIF()
ENDIF()

IF (ASSIMP_ASAN)
  MESSAGE(STATUS "AddressSanitizer enabled")
  SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fsanitize=address")
  SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fsanitize=address")
ENDIF()

IF (ASSIMP_UBSAN)
  MESSAGE(STATUS "Undefined Behavior sanitizer enabled")
  SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fsanitize=undefined,shift,shift-exponent,integer-divide-by-zero,unreachable,vla-bound,null,return,signed-integer-overflow,bounds,float-divide-by-zero,float-cast-overflow,nonnull-attribute,returns-nonnull-attribute,bool,enum,vptr,pointer-overflow,builtin -fno-sanitize-recover=all")
  SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fsanitize=undefined,shift,shift-exponent,integer-divide-by-zero,unreachable,vla-bound,null,return,signed-integer-overflow,bounds,float-divide-by-zero,float-cast-overflow,nonnull-attribute,returns-nonnull-attribute,bool,enum,vptr,pointer-overflow,builtin -fno-sanitize-recover=all")
ENDIF()

INCLUDE (FindPkgMacros)
INCLUDE (PrecompiledHeader)

# Set Assimp project output directory variables.
SET(ASSIMP_RUNTIME_OUTPUT_DIRECTORY "${CMAKE_CURRENT_BINARY_DIR}/bin" CACHE STRING "Path for runtime output files")
SET(ASSIMP_LIBRARY_OUTPUT_DIRECTORY "${CMAKE_CURRENT_BINARY_DIR}/bin" CACHE STRING "Path for library output files")
SET(ASSIMP_ARCHIVE_OUTPUT_DIRECTORY "${CMAKE_CURRENT_BINARY_DIR}/lib" CACHE STRING "Path for archive output files")

# Macro used to set the output directories of a target to the
# respective Assimp output directories.
MACRO(TARGET_USE_COMMON_OUTPUT_DIRECTORY target)
  set_target_properties(${target} PROPERTIES
      RUNTIME_OUTPUT_DIRECTORY ${ASSIMP_RUNTIME_OUTPUT_DIRECTORY}
      LIBRARY_OUTPUT_DIRECTORY ${ASSIMP_LIBRARY_OUTPUT_DIRECTORY}
      ARCHIVE_OUTPUT_DIRECTORY ${ASSIMP_ARCHIVE_OUTPUT_DIRECTORY}
  )
ENDMACRO()

get_cmake_property(is_multi_config GENERATOR_IS_MULTI_CONFIG)

IF (ASSIMP_INJECT_DEBUG_POSTFIX AND (is_multi_config OR CMAKE_BUILD_TYPE STREQUAL "Debug"))
  SET(CMAKE_DEBUG_POSTFIX "d" CACHE STRING "Debug Postfix for lib, samples and tools")
ELSE()
  SET(CMAKE_DEBUG_POSTFIX "" CACHE STRING "Debug Postfix for lib, samples and tools")
ENDIF()

# Only generate this target if no higher-level project already has
IF (NOT TARGET uninstall AND ASSIMP_INSTALL)
  # add make uninstall capability
  CONFIGURE_FILE("${CMAKE_CURRENT_SOURCE_DIR}/cmake-modules/cmake_uninstall.cmake.in" "${CMAKE_CURRENT_BINARY_DIR}/cmake_uninstall.cmake" IMMEDIATE @ONLY)
  ADD_CUSTOM_TARGET(uninstall "${CMAKE_COMMAND}" -P "${CMAKE_CURRENT_BINARY_DIR}/cmake_uninstall.cmake")
ENDIF()

# cmake configuration files
if(${BUILD_SHARED_LIBS})
  set(BUILD_LIB_TYPE SHARED)
else()
  set(BUILD_LIB_TYPE STATIC)
endif()

set(GENERATED_DIR "${CMAKE_CURRENT_BINARY_DIR}/generated")

set(INCLUDE_INSTALL_DIR "include")

IF( ASSIMP_BUILD_DOCS )
  ADD_SUBDIRECTORY(doc)
ENDIF()

    INCLUDE(CheckIncludeFile)
    INCLUDE(CheckTypeSize)
    INCLUDE(CheckFunctionExists)

    # Explicitly turn off ASM686 and AMD64 cmake options.
    # The AMD64 option causes a build failure on MSVC and the ASM builds seem to have problems:
    #		https://github.com/madler/zlib/issues/41#issuecomment-125848075
    # Also prevents these options from "polluting" the cmake options if assimp is being
    # included as a submodule.
    set( ASM686 FALSE CACHE INTERNAL "Override ZLIB flag to turn off assembly" FORCE )
    set( AMD64 FALSE CACHE INTERNAL "Override ZLIB flag to turn off assembly" FORCE )

    # compile from sources
    ADD_SUBDIRECTORY(contrib/zlib)
    SET(ZLIB_FOUND 1)
    SET(ZLIB_LIBRARIES zlibstatic)
    SET(ZLIB_INCLUDE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/contrib/zlib ${CMAKE_CURRENT_BINARY_DIR}/contrib/zlib)
    # need to ensure we don't link with system zlib or minizip as well.
    SET(ASSIMP_BUILD_MINIZIP 1)

IF( NOT IOS )
  IF( NOT ASSIMP_BUILD_MINIZIP )
    use_pkgconfig(UNZIP minizip)
  ENDIF()
ELSE ()
  IF( NOT BUILD_SHARED_LIBS )
    IF( NOT ASSIMP_BUILD_MINIZIP )
      use_pkgconfig(UNZIP minizip)
    ENDIF()
  ENDIF ()
ENDIF ()

IF ( ASSIMP_NO_EXPORT )
  ADD_DEFINITIONS( -DASSIMP_BUILD_NO_EXPORT)
  MESSAGE( STATUS "Build an import-only version of Assimp." )
ENDIF()

SET ( ASSIMP_BUILD_ARCHITECTURE "" CACHE STRING
  "describe the current architecture."
)
IF( ASSIMP_BUILD_ARCHITECTURE STREQUAL "")
ELSE()
  ADD_DEFINITIONS ( -D'ASSIMP_BUILD_ARCHITECTURE="${ASSIMP_BUILD_ARCHITECTURE}"' )
ENDIF()

# ${CMAKE_GENERATOR}
SET ( ASSIMP_BUILD_COMPILER "" CACHE STRING
  "describe the current compiler."
)
IF( ASSIMP_BUILD_COMPILER STREQUAL "")
ELSE()
  ADD_DEFINITIONS ( -D'ASSIMP_BUILD_COMPILER="${ASSIMP_BUILD_COMPILER}"' )
ENDIF()

MARK_AS_ADVANCED ( ASSIMP_BUILD_ARCHITECTURE ASSIMP_BUILD_COMPILER )

# Draco requires cmake 3.12
IF (DEFINED CMAKE_VERSION AND "${CMAKE_VERSION}" VERSION_LESS "3.12")
  message(NOTICE "draco requires cmake 3.12 or newer, cmake is ${CMAKE_VERSION} . Draco is disabled")
  SET ( ASSIMP_BUILD_DRACO OFF CACHE BOOL "Disabled: Draco requires newer cmake" FORCE )
ELSE()
  OPTION ( ASSIMP_BUILD_DRACO "If the Draco libraries are to be built. Primarily for glTF" OFF )
  IF ( ASSIMP_BUILD_DRACO )
    # Primarily for glTF v2
    # Enable Draco glTF feature set
    set(DRACO_GLTF ON CACHE BOOL "" FORCE)
    # Disable unnecessary or omitted components
    set(DRACO_JS_GLUE OFF CACHE BOOL "" FORCE)
    set(DRACO_WASM OFF CACHE BOOL "" FORCE)
    set(DRACO_MAYA_PLUGIN OFF CACHE BOOL "" FORCE)
    set(DRACO_UNITY_PLUGIN OFF CACHE BOOL "" FORCE)
    set(DRACO_TESTS OFF CACHE BOOL "" FORCE)

    IF(ASSIMP_HUNTER_ENABLED)
      hunter_add_package(draco)
      find_package(draco CONFIG REQUIRED)
      set(draco_LIBRARIES draco::draco)
    ELSE()
      # Draco 1.4.1 has many warnings and will not build with /WX or -Werror
      # See https://github.com/google/draco/issues/672
      # and https://github.com/google/draco/issues/673
      IF(MSVC)
        set(DRACO_CXX_FLAGS "/W0")
      ELSE()
        list(APPEND DRACO_CXX_FLAGS
          "-Wno-bool-compare"
          "-Wno-comment"
          "-Wno-maybe-uninitialized"
          "-Wno-sign-compare"
          "-Wno-unused-local-typedefs"
        )
        # Draco 1.4.1 does not explicitly export any symbols under GCC/clang
        list(APPEND DRACO_CXX_FLAGS
          "-fvisibility=default"
        )
      ENDIF()

      # Don't build or install all of Draco by default
      ADD_SUBDIRECTORY( "contrib/draco" EXCLUDE_FROM_ALL )

      if(MSVC OR WIN32)
        set(draco_LIBRARIES "draco")
      else()
        if(BUILD_SHARED_LIBS)
          set(draco_LIBRARIES "draco_shared")
        else()
          set(draco_LIBRARIES "draco_static")
        endif()
      endif()

      # Don't build the draco command-line tools by default
      set_target_properties(draco_encoder draco_decoder PROPERTIES
        EXCLUDE_FROM_ALL TRUE
        EXCLUDE_FROM_DEFAULT_BUILD TRUE
        )

      # Do build the draco shared library
      set_target_properties(${draco_LIBRARIES} PROPERTIES
        EXCLUDE_FROM_ALL FALSE
        EXCLUDE_FROM_DEFAULT_BUILD FALSE
        )

      TARGET_USE_COMMON_OUTPUT_DIRECTORY(${draco_LIBRARIES})
      TARGET_USE_COMMON_OUTPUT_DIRECTORY(draco_encoder)
      TARGET_USE_COMMON_OUTPUT_DIRECTORY(draco_decoder)

      set(draco_INCLUDE_DIRS "${CMAKE_CURRENT_SOURCE_DIR}/contrib/draco/src")

      # This is probably wrong
      INSTALL( TARGETS ${draco_LIBRARIES}
        EXPORT "${TARGETS_EXPORT_NAME}"
        LIBRARY DESTINATION ${ASSIMP_LIB_INSTALL_DIR}
        ARCHIVE DESTINATION ${ASSIMP_LIB_INSTALL_DIR}
        RUNTIME DESTINATION ${ASSIMP_BIN_INSTALL_DIR}
        FRAMEWORK DESTINATION ${ASSIMP_LIB_INSTALL_DIR}
        COMPONENT ${LIBASSIMP_COMPONENT}
        INCLUDES DESTINATION include
    )

    ENDIF()
  ENDIF()
ENDIF()

# Main assimp code
ADD_SUBDIRECTORY( code/ )